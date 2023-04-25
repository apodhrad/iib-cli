package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
	TestLogger    *log.Logger
)

func init() {
	logFile := "/tmp/logs.txt"
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	TestLogger = log.New(file, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile)
}

const GRPC_NAME string = "iib_registry_server"
const GRPC_HOST string = "localhost"
const GRPC_PORT string = "50051"
const GRPC_SERVER string = GRPC_HOST + ":" + GRPC_PORT

type GrpcArg struct {
	data   string
	api    string
	method string
}

func GrpcArgApi(api string) GrpcArg {
	return GrpcArg{api: api}
}

func GrpcArgMethod(method string) GrpcArg {
	return GrpcArg{method: method}
}

func GrpcArgMethodWithData(method string, data string) GrpcArg {
	return GrpcArg{method: method, data: data}
}

func GrpcArgToCmdArgs(grpcArg GrpcArg) ([]string, error) {
	var cmdArgs []string = []string{"-plaintext"}
	if grpcArg.data != "" {
		cmdArgs = append(cmdArgs, "-d", grpcArg.data)
	}
	cmdArgs = append(cmdArgs, GRPC_SERVER)
	if grpcArg.api != "" {
		apiArgs := strings.Split(grpcArg.api, " ")
		cmdArgs = append(cmdArgs, apiArgs...)
	} else if grpcArg.method != "" {
		cmdArgs = append(cmdArgs, grpcArg.method)
	} else {
		return cmdArgs, errors.New("No api or method is defined")
	}
	return cmdArgs, nil
}

func GrpcStart() {
	iib := os.Getenv("IIB")
	if iib == "" {
		ErrorLogger.Panicln("Specify index image via envvar IIB or via command set iib!")
	}

	// make sure there is no running container before starting a new one
	InfoLogger.Println("Make sure there is no running container")
	GrpcStop()

	// initialize docker client
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		ErrorLogger.Panicln(err)
	}
	defer cli.Close()

	// pull the index image bundle (iib)
	InfoLogger.Println("Pull image " + iib)
	out, err := cli.ImagePull(ctx, iib, types.ImagePullOptions{})
	if err != nil {
		ErrorLogger.Panicln(err)
	}
	defer out.Close()

	// now, we can start the new container
	InfoLogger.Println("Crate a new container")
	resp, err := cli.ContainerCreate(
		ctx,
		&container.Config{
			Image:        iib,
			ExposedPorts: nat.PortSet{"50051": struct{}{}},
		},
		&container.HostConfig{
			PortBindings: nat.PortMap{
				"50051": []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: GRPC_PORT,
					},
				},
			},
		},
		nil,
		nil,
		GRPC_NAME)
	if err != nil {
		ErrorLogger.Panicln(err.Error())
	}

	InfoLogger.Println("Start the new container")
	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		ErrorLogger.Panicln(err)
	}

	// the container should be created, now wait for its running state
	InfoLogger.Println("Wait for its running state")
	grpcContainer, err := waitForState(GRPC_NAME, "running")
	if err != nil {
		ErrorLogger.Panicln(err)
	}

	// the container should be running, now wait for its readiness
	InfoLogger.Println("Wait for its readiness")
	err = waitForResponse()
	if err != nil {
		ErrorLogger.Panicln(err)
	}
	InfoLogger.Println("The container is up and running " + containerToString(grpcContainer))
}

func GrpcStop() {
	// check if there already is a container
	grpcContainer := getContainerWithName(GRPC_NAME)
	if grpcContainer == nil {
		// if not then there is nothing to stop
		return
	}

	// a container exists, so make sure it is in a proper state before its removal
	InfoLogger.Println("Wait for running state in " + containerToString(grpcContainer))
	grpcContainer, err := waitForState(GRPC_NAME, "running")
	if err != nil {
		ErrorLogger.Panicln(err)
	}

	// initialize docker client
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		ErrorLogger.Panicln(err)
	}
	defer cli.Close()

	// once the conatainer is in running state, we can stop it
	InfoLogger.Println("Stop the container " + containerToString(grpcContainer))
	err = cli.ContainerStop(ctx, grpcContainer.ID, container.StopOptions{})
	if err != nil {
		ErrorLogger.Panicln(err)
	}
	InfoLogger.Println("Wait for exited state in " + containerToString(grpcContainer))
	grpcContainer, err = waitForState(GRPC_NAME, "exited")
	if err != nil {
		ErrorLogger.Panicln(err)
	}

	// once the conatainer is in exited state, we can remove it
	InfoLogger.Println("Remove the container " + containerToString(grpcContainer))
	err = cli.ContainerRemove(ctx, grpcContainer.ID, types.ContainerRemoveOptions{})
	if err != nil {
		ErrorLogger.Panicln(err.Error())
	}
	InfoLogger.Println("Wait until the container is gone")
	grpcContainer, err = waitForState(GRPC_NAME, "")
	if err != nil {
		ErrorLogger.Panicln(err)
	}
}

func GrpcStatus() (string, error) {
	grpcContainer := getContainerWithName(GRPC_NAME)
	if grpcContainer != nil {
		return grpcContainer.Status, nil
	}
	return "", nil
}

func containerToString(c *types.Container) string {
	if c == nil {
		return "Container[nil]"
	}
	return fmt.Sprintf("Container[Name: %s, State: %s, Status: %s]", c.Names[0], c.State, c.Status)
}

func getContainerWithName(name string) *types.Container {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		panic(err)
	}

	for _, container := range containers {
		for _, name := range container.Names {
			if name == "/"+GRPC_NAME {
				return &container
			}
		}
	}

	return nil
}

func waitForState(name string, state string) (*types.Container, error) {
	var c *types.Container

	for i := 0; i < 10; i++ {
		c = getContainerWithName(name)
		if c != nil && c.State == "recovered" {
			panic("State is recovered. " + containerToString(c))
		}
		if state != "" {
			if c != nil && c.State == state {
				return c, nil
			}
		} else {
			if c == nil {
				return nil, nil
			}
		}
		time.Sleep(1 * time.Second)
	}

	if state != "" {
		return nil, fmt.Errorf("%s is still not in a state '%s'", containerToString(c), state)
	} else {
		return nil, fmt.Errorf("%s is still present", containerToString(c))
	}
}

func waitForResponse() error {
	var err error
	for i := 0; i < 10; i++ {
		out, err := GrpcExec(GrpcArgApi("list"))
		if err == nil && out != "" {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return err
}

func GrpcExec(grpcArg GrpcArg) (string, error) {
	cmdArgs, err := GrpcArgToCmdArgs(grpcArg)
	if err != nil {
		return "", err
	}

	var errOut bytes.Buffer
	cmd := exec.Command("grpcurl", cmdArgs...)
	cmd.Stderr = &errOut
	out, err := cmd.Output()
	if err != nil {
		err = errors.New(errOut.String() + err.Error())
	}
	return string(out), err
}
