package grpc

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

var (
	DebugLogger   *log.Logger
	InfoLogger    *log.Logger
	WarningLogger *log.Logger
	ErrorLogger   *log.Logger
)

func init() {
	logFile := "/tmp/logs.txt"
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	DebugLogger = log.New(file, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
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
	// check is iib was specified
	iib := os.Getenv("IIB")
	if iib == "" {
		ErrorLogger.Panicln("Specify index image via envvar IIB or via command set iib!")
	}

	// make sure there is no running container before starting a new one
	InfoLogger.Println("Make sure there is no running container")
	GrpcStop()

	// pull the index image bundle (iib)
	InfoLogger.Printf("Pull image %s\n", iib)
	err := DockerPullImage(iib)
	if err != nil {
		ErrorLogger.Panicln(err)
	}

	// now, we can start the new container
	InfoLogger.Println("Start grpc server on localhost:" + GRPC_PORT)
	id, err := DockerStartContainer(GRPC_NAME, iib, GRPC_PORT+":"+GRPC_PORT)
	if err != nil {
		ErrorLogger.Panicln(err)
	}
	InfoLogger.Printf("Container wit ID %s was sucesfully started", id)

	// the container should be running, now wait for its readiness
	InfoLogger.Println("Wait for its readiness")
	err = waitForResponse()
	if err != nil {
		ErrorLogger.Panicln(err)
	}
	InfoLogger.Println("The grpc server is up and running on localhost:" + GRPC_PORT)
}

func GrpcStop() {
	InfoLogger.Println("Stop grpc server")
	err := DockerStopContainer(GRPC_NAME)
	if err != nil {
		ErrorLogger.Panicln(err)
	}
	InfoLogger.Println("The grpc server is stopped")
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
