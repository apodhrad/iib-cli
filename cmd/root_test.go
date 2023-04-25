package cmd

import (
	"log"
	"os"
	"testing"

	"github.com/apodhrad/iib-cli/grpc"
)

var (
	TestLogger *log.Logger
)

func init() {
	logFile := "/tmp/logs.txt"
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	TestLogger = log.New(file, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile)
}
func setTestIIB(t *testing.T) {
	TestLogger.Println("Set IIB for test " + t.Name())
	os.Setenv("IIB", "quay.io/apodhrad/iib-test:v0.0.1")
}

func startTestGrpc(t *testing.T) {
	setTestIIB(t)
	grpc.GrpcStart()
}

func stopTestGrpc(t *testing.T) {
	os.Unsetenv("IIB")
	grpc.GrpcStop()
}
