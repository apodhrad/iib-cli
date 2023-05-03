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
]`

const EXPECTED_PACKAGES_TEXT string = `PACKAGE_NAME    
prometheus      
redis-operator  
`

func TestPackagesCmdJson(t *testing.T) {
	setTestIIB(t)
	defer stopTestGrpc(t)

	funcArgs := PackagesCmdArgs{Output: "json"}
	out, err := packagesCmdFunc(funcArgs)
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_PACKAGES_JSON, out)
}

func TestPackagesCmdText(t *testing.T) {
	setTestIIB(t)
	defer stopTestGrpc(t)

	funcArgs := PackagesCmdArgs{Output: "text"}
	out, err := packagesCmdFunc(funcArgs)
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_PACKAGES_TEXT, out)
}
