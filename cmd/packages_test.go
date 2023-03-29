package cmd

import (
	"os"
	"testing"

	"github.com/apodhrad/iib-cli/utils"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {

}

func TestPackagesRunE(t *testing.T) {
	utils.GrpcStopSafely()

	os.Setenv("IIB", "quay.io/apodhrad/iib-test:v0.0.1")
	utils.GrpcStartSafely()

	err := packagesRunE(nil, []string{})
	assert.Nil(t, err)

	utils.GrpcStopSafely()
}

const EXPECTED_PACKAGES_GRPC_OUTPUT string = `{
	"name": "prometheus"
}
{
	"name": "redis-operator"
}
`

var EXPECTED_PACKAGES []Package = []Package{{Name: "prometheus"}, {Name: "redis-operator"}}

func TestGetPackages(t *testing.T) {
	packages, err := getPackages(EXPECTED_PACKAGES_GRPC_OUTPUT)
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_PACKAGES, packages)
}
