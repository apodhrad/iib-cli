package cmd

import (
	"os"
	"testing"

	"github.com/apodhrad/iib-cli/grpc"
)

func setTestIIB(t *testing.T) {
	grpc.TestLogger.Println("Set IIB for test " + t.Name())
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
