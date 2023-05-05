package grpc

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
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
	GrpcStart()
}

func stopTestGrpc(t *testing.T) {
	os.Unsetenv("IIB")
	GrpcStop()
}

func TestGrpcConstants(t *testing.T) {
	assert.Equal(t, "localhost:50051", GRPC_SERVER)
}
