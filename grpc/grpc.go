package grpc

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/apodhrad/iib-cli/logging"
)

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

func GrpcStart() string {
	// check is iib was specified
	iib := os.Getenv("IIB")
	if iib == "" {
		handlePanic(fmt.Errorf("Specify index image via envvar IIB or via command set iib!"))
	}

	// make sure there is no running container before starting a new one
	logging.INFO().Printf("Make sure there is no running container")
	GrpcStop()

	// pull the index image bundle (iib)
	logging.INFO().Printf("Pull image %s\n", iib)
	err := DockerPullImage(iib)
	if err != nil {
		handlePanic(err)
	}

	// now, we can start the new container
	logging.INFO().Printf("Start grpc server on localhost:" + GRPC_PORT)
	id, err := DockerStartContainer(GRPC_NAME, iib, GRPC_PORT+":"+GRPC_PORT)
	if err != nil {
		handlePanic(err)
	}
	logging.INFO().Printf("Container wit ID %s was sucesfully started", id)

	address := GRPC_HOST + ":" + GRPC_PORT
	// the container should be running, now wait for its readiness
	logging.INFO().Printf("Wait for its readiness")
	err = waitForReadiness(address)
	if err != nil {
		handlePanic(err)
	}
	logging.INFO().Printf("The grpc server is up and running on %s", address)
	return address
}

func GrpcStop() {
	logging.INFO().Printf("Stop grpc server")
	err := DockerStopContainer(GRPC_NAME)
	if err != nil {
		handlePanic(err)
	}
	logging.INFO().Printf("The grpc server is stopped")
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

func waitForReadiness(address string) error {
	var err error
	for i := 0; i < 10; i++ {
		client, err := NewClient(address)
		defer client.Close()
		if client != nil && err == nil {
			isReady, _ := client.HealthCheck()
			if isReady {
				return nil
			}
		}
		time.Sleep(1 * time.Second)
	}
	return err
}

func handlePanic(err error) {
	logging.ERROR().Println(err)
	panic(err)
}
