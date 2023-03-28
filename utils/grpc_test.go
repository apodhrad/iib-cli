package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXPECTED_API_LIST string = `api.Registry
grpc.health.v1.Health
grpc.reflection.v1alpha.ServerReflection
`

const EXPECTED_API_LIST_REGISTRY string = `api.Registry.GetBundle
api.Registry.GetBundleForChannel
api.Registry.GetBundleThatReplaces
api.Registry.GetChannelEntriesThatProvide
api.Registry.GetChannelEntriesThatReplace
api.Registry.GetDefaultBundleThatProvides
api.Registry.GetLatestChannelEntriesThatProvide
api.Registry.GetPackage
api.Registry.ListBundles
api.Registry.ListPackages
`

const EXPECTED_API_DESCRIPTION_GETPKGREQ string = `api.GetPackageRequest is a message:
message GetPackageRequest {
  string name = 1;
}
`

func setIIB() {
	os.Setenv("IIB", "quay.io/apodhrad/iib-test:v0.0.1")
}

func start() {
	setIIB()
	GrpcStartSafely()
}

func clean() {
	os.Unsetenv("IIB")
	GrpcStop()
}

func TestGrpcConstants(t *testing.T) {
	assert.Equal(t, "localhost:50051", GRPC_SERVER)
}

func TestGrpcStartStop(t *testing.T) {
	setIIB()

	err := GrpcStartSafely()
	assert.Nil(t, err)
	status, err := GrpcStatus()
	assert.Regexp(t, "^Up", status)

	err = GrpcStop()
	assert.Nil(t, err)
	status, err = GrpcStatus()
	assert.Empty(t, status)

	clean()
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

func TestGrpcApiList(t *testing.T) {
	start()

	stdOut, _ := GrpcExec(GrpcArg{api: "list"})
	assert.Equal(t, EXPECTED_API_LIST, stdOut)

	stdOut, _ = GrpcExec(GrpcArg{api: "list api.Registry"})
	assert.Equal(t, EXPECTED_API_LIST_REGISTRY, stdOut)

	clean()
}

func TestGrpcApiDescribe(t *testing.T) {
	start()

	stdOut, _ := GrpcExec(GrpcArg{api: "describe api.GetPackageRequest"})
	assert.Equal(t, EXPECTED_API_DESCRIPTION_GETPKGREQ, stdOut)

	clean()
}

const EXPECTED_REGISTRY_PACKAGES string = `{
  "name": "prometheus"
}
{
  "name": "redis-operator"
}
`

func TestGrpcRegistryListPackages(t *testing.T) {
	start()

	stdOut, _ := GrpcExec(GrpcArg{method: "api.Registry/ListPackages"})
	assert.Equal(t, EXPECTED_REGISTRY_PACKAGES, stdOut)

	clean()
}
