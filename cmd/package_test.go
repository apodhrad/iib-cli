package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const EXPECTED_PACKAGE_REDIS_GRPC_OUTPUT string = `{
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

// func TestGetPackage(t *testing.T) {
// 	setup()

// 	err := packageRunE(nil, []string{"redis-operator"})
// 	assert.Nil(t, err)

// 	teardown()
// }

func TestPackageCmdGrpc(t *testing.T) {
	setup()

	out, err := packageCmdGrpc("redis-operator")
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_PACKAGE_REDIS_GRPC_OUTPUT, out)

	teardown()
}

var EXPETED_PACKAGE_REDIS Package = Package{Name: "redis-operator", Channels: []Channel{{Name: "preview", CsvName: "redis-operator.v0.4.0"}, {Name: "stable", CsvName: "redis-operator.v0.13.0"}}, DefaultChannelName: "stable"}

func TestPackageCmdUnmarshal(t *testing.T) {
	pkg, err := packageCmdUnmarshal(EXPECTED_PACKAGE_REDIS_GRPC_OUTPUT)
	assert.Nil(t, err)
	assert.Equal(t, EXPETED_PACKAGE_REDIS, pkg)
}
