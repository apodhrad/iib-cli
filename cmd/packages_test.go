package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXPECTED_PACKAGES_JSON string = `[
  {
    "name": "prometheus"
  },
  {
    "name": "redis-operator"
  }
]
`

const EXPECTED_PACKAGES_TABLE string = `PACKAGE_NAME    
prometheus      
redis-operator  
`

func TestCmdGetPackages(t *testing.T) {
	out, err := testCmd(t, "get", "packages")

	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_PACKAGES_TABLE, out)
}

func TestCmdGetPackagesJson(t *testing.T) {
	out, err := testCmd(t, "get", "packages", "-o", "json")

	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_PACKAGES_JSON, string(out))
}
