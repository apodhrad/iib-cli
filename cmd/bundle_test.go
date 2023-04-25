package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBundleCmdGprc(t *testing.T) {
	setTestIIB(t)
	defer stopTestGrpc(t)

	out, err := bundleCmdGrpc("redis-operator.v0.8.0", "redis-operator", "stable")
	assert.Nil(t, err)
	expected := readFile("bundle.json")
	assert.Equal(t, expected, out)
}

var EXPECTED_BUNDLE = Bundle{
	CsvName:     "redis-operator.v0.8.0",
	PackageName: "redis-operator",
	ChannelName: "stable",
	BundlePath:  "quay.io/openshift-community-operators/redis-operator:v0.8.0",
	Version:     "0.8.0",
}

func TestBundleCmdUnmarshal(t *testing.T) {
	input := readFile("bundle.json")
	bundle, err := bundleCmdUnmarshal(input)
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_BUNDLE, bundle)
}

const EXPECTED_BUNDLE_TEXT_OUTPUT string = `CSV                    PACKAGE         CHANNEL  
redis-operator.v0.8.0  redis-operator  stable   
`

func TestBundleCmdToText(t *testing.T) {
	out, err := bundleCmdToText(EXPECTED_BUNDLE)
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_BUNDLE_TEXT_OUTPUT, out)
}

const EXPECTED_BUNDLE_JSON_OUTPUT string = `{"csvName":"redis-operator.v0.8.0","packageName":"redis-operator","channelname":"stable","bundlePath":"quay.io/openshift-community-operators/redis-operator:v0.8.0","version":"0.8.0"}`

func TestBundleCmdToJson(t *testing.T) {
	out, err := bundleCmdToJson(EXPECTED_BUNDLE)
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_BUNDLE_JSON_OUTPUT, out)
}
