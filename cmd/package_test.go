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
}`

const EXPECTED_PACKAGE_REDIS_TABLE string = `PACKAGE_NAME    CHANNEL  CSV                     DEFAULT  
redis-operator  preview  redis-operator.v0.4.0            
redis-operator  stable   redis-operator.v0.13.0  true     
`

func TestPackageCmdJson(t *testing.T) {
	setTestIIB(t)
	defer stopTestGrpc(t)

	funcArgs := PackageCmdArgs{Name: "redis-operator", Output: "json"}
	out, err := packageCmdFunc(funcArgs)
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_PACKAGE_REDIS_JSON, out)
}

func TestPackageCmdText(t *testing.T) {
	setTestIIB(t)
	defer stopTestGrpc(t)

	funcArgs := PackageCmdArgs{Name: "redis-operator", Output: "text"}
	out, err := packageCmdFunc(funcArgs)
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_PACKAGE_REDIS_TABLE, out)
}
