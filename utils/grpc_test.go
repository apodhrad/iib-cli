package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup() {
	GrpcStopSafely()
	os.Setenv("IIB", "quay.io/apodhrad/iib-test:v0.0.1")
	GrpcStartSafely()
}

func teardown() {
	GrpcStopSafely()
	os.Unsetenv("IIB")
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

const EXPECTED_API_DESCRIPTION_GETPKGREQ string = `api.GetPackageRequest is a message:
message GetPackageRequest {
  string name = 1;
}
`

func TestGrpcStartStopWithRequest(t *testing.T) {
	var err error
	var out string
	var status string

	os.Setenv("IIB", "quay.io/apodhrad/iib-test:v0.0.1")

	err = GrpcStartSafely()
	assert.Nil(t, err)
	status, err = GrpcStatus()
	assert.Nil(t, err)
	assert.Regexp(t, "^Up", status)

	out, err = GrpcExec(GrpcArg{api: "describe api.GetPackageRequest"})
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_API_DESCRIPTION_GETPKGREQ, out)

	err = GrpcStopSafely()
	assert.Nil(t, err)
	status, err = GrpcStatus()
	assert.Nil(t, err)
	assert.Empty(t, status)

	teardown()
}
