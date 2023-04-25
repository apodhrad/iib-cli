package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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

func TestGrcpArgToCmdArgs(t *testing.T) {
	var cmdArgs []string
	var err error

	cmdArgs, err = GrpcArgToCmdArgs(GrpcArg{api: "list"})
	assert.Nil(t, err)
	assert.Equal(t, []string{"-plaintext", GRPC_SERVER, "list"}, cmdArgs)

	cmdArgs, err = GrpcArgToCmdArgs(GrpcArg{api: "list API"})
	assert.Nil(t, err)
	assert.Equal(t, []string{"-plaintext", GRPC_SERVER, "list", "API"}, cmdArgs)

	cmdArgs, err = GrpcArgToCmdArgs(GrpcArg{api: "describe API"})
	assert.Nil(t, err)
	assert.Equal(t, []string{"-plaintext", GRPC_SERVER, "describe", "API"}, cmdArgs)

	cmdArgs, err = GrpcArgToCmdArgs(GrpcArg{method: "Service/Method"})
	assert.Nil(t, err)
	assert.Equal(t, []string{"-plaintext", GRPC_SERVER, "Service/Method"}, cmdArgs)

	var data string = "{\"pkgName\":\"my-package\"}"
	cmdArgs, err = GrpcArgToCmdArgs(GrpcArg{data: data, method: "Service/Method"})
	assert.Nil(t, err)
	assert.Equal(t, []string{"-plaintext", "-d", data, GRPC_SERVER, "Service/Method"}, cmdArgs)

	cmdArgs, err = GrpcArgToCmdArgs(GrpcArg{})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "No api or method is defined")

	cmdArgs, err = GrpcArgToCmdArgs(GrpcArg{data: data})
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "No api or method is defined")
}
