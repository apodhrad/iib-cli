package utils

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
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

func GrpcStart() error {
	iib := os.Getenv("IIB")
	if iib == "" {
		return errors.New("Specify index image via envvar IIB or via command set iib")
	}

	cmd := exec.Command("podman", "run", "-d", "--name", GRPC_NAME, "-p", GRPC_PORT+":"+GRPC_PORT, iib)
	err := cmd.Run()
	return err
}

func GrpcStartSafely() error {
	var err error
	var out string
	var status string

	status, err = GrpcStatus()
	if err != nil {
		return err
	}
	regex := regexp.MustCompile("^Up")
	if regex.MatchString(status) {
		// ok, the server is already started
		return nil
	}

	err = GrpcStopSafely()
	if err != nil {
		return err
	}

	err = GrpcStart()
	if err != nil {
		return err
	}

	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		status, err = GrpcStatus()
		if err != nil {
			return err
		}
		if regex.MatchString(status) {
			out, err = GrpcExec(GrpcArgApi("list"))
			if err == nil && out != "" {
				// ok, the server is up and responding
				return nil
			}
		}
	}

	return errors.New("Server was not started properly. Status: " + status)
}

func GrpcStop() (string, error) {
	cmd := exec.Command("podman", "rm", "-f", "-i", GRPC_NAME)
	out, err := cmd.Output()
	return string(out), err
}

func GrpcStopSafely() error {
	var err error
	var out string
	var status string

	out, err = GrpcStop()
	if err != nil {
		return err
	}
	if out != "" {
		for i := 0; i < 10; i++ {
			time.Sleep(2 * time.Second)
			status, err = GrpcStatus()
			if err != nil {
				return err
			}
			if status == "" {
				return nil
			}
		}
		return errors.New("Server was not stopped properly. Status: " + status)
	}

	return nil
}

func GrpcStatus() (string, error) {
	cmd := exec.Command("podman", "ps", "-a", "--format", "{{.Status}}", "-f", "name="+GRPC_NAME)
	out, err := cmd.Output()
	var status string = string(out)
	status = strings.Replace(status, "\n", "", -1)
	return status, err
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
