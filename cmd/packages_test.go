package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXPECTED_PACKAGES_GRPC_OUTPUT string = `{
  "name": "prometheus"
}
{
  "name": "redis-operator"
}
`

func TestPackagesCmdGrpc(t *testing.T) {
	setTestIIB(t)
	defer stopTestGrpc(t)

	out, err := packagesCmdGrpc()
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_PACKAGES_GRPC_OUTPUT, out)
}

var EXPECTED_PACKAGES []Package = []Package{{Name: "prometheus"}, {Name: "redis-operator"}}

func TestPackagesCmdUnmarshal(t *testing.T) {
	packages, err := packagesCmdUnmarshal(EXPECTED_PACKAGES_GRPC_OUTPUT)
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_PACKAGES, packages)
}

const EXPECTED_PACKAGES_TEST_OUTPUT string = `PACKAGE         
prometheus      
redis-operator  
`

func TestPackagesCmdToText(t *testing.T) {
	out, err := packagesCmdToText(EXPECTED_PACKAGES)
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_PACKAGES_TEST_OUTPUT, out)
}
