package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXPECTED_PACKAGE_REDIS_JSON string = `{
  "name": "redis-operator",
  "channels": [
    {
      "name": "preview",
      "csvName": "redis-operator.v0.4.0"
    },
    {
      "name": "stable",
      "csvName": "redis-operator.v0.13.0"
    }
  ],
  "defaultChannelName": "stable"
}
`

const EXPECTED_PACKAGE_REDIS_TABLE string = `PACKAGE_NAME    CHANNEL  CSV                     DEFAULT  
redis-operator  preview  redis-operator.v0.4.0            
redis-operator  stable   redis-operator.v0.13.0  true     
`

func TestCmdGetPackage(t *testing.T) {
	out, err := testCmd(t, "get", "package", "redis-operator")

	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_PACKAGE_REDIS_TABLE, out)
}

func TestCmdGetPackageJson(t *testing.T) {
	out, err := testCmd(t, "get", "package", "redis-operator", "-o", "json")

	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_PACKAGE_REDIS_JSON, out)
}

func TestCmdGetPackageNone(t *testing.T) {
	out, err := testCmd(t, "get", "package")

	assert.NotNil(t, err)
	assert.Equal(t, "", out)
}
