package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBundlesCmdGprc(t *testing.T) {
	setTestIIB(t)
	defer stopTestGrpc(t)

	out, err := bundlesCmdGrpc()
	assert.Nil(t, err)
	expected := readFile("bundles.json")
	assert.Equal(t, expected, out)
}

var EXPECTED_BUNDLE_PROMETHEUS Bundle = Bundle{
	CsvName:     "prometheusoperator.0.47.0",
	PackageName: "prometheus",
	ChannelName: "beta",
	BundlePath:  "quay.io/openshift-community-operators/prometheus:v0.47.0",
	Version:     "0.47.0",
	Replaces:    "prometheusoperator.0.37.0",
}

var EXPECTED_BUNDLE_REDIS Bundle = Bundle{
	CsvName:     "redis-operator.v0.8.0",
	PackageName: "redis-operator",
	ChannelName: "stable",
	BundlePath:  "quay.io/openshift-community-operators/redis-operator:v0.8.0",
	Version:     "0.8.0",
	Replaces:    "redis-operator.v0.6.0",
}

func TestBundlesCmdUnmarshal(t *testing.T) {
	input := readFile("bundles.json")
	bundles, err := bundlesCmdUnmarshal(input)
	assert.Nil(t, err)
	// assert.Equal(t, EXPECTED_BUNDLE, bundle)
	prometheusFound := false
	redisFound := false
	for _, bundle := range bundles {
		if bundle == EXPECTED_BUNDLE_PROMETHEUS {
			prometheusFound = true
		}
		if bundle == EXPECTED_BUNDLE_REDIS {
			redisFound = true
		}
	}
	assert.True(t, prometheusFound, "Bundle 'prometheusoperator.0.47.0' not found")
	assert.True(t, redisFound, "Bundle 'redis-operator.v0.8.0' not found")
}

var EXPECTED_BUNDLES []Bundle = []Bundle{EXPECTED_BUNDLE_PROMETHEUS, EXPECTED_BUNDLE_REDIS}

const EXPECTED_BUNDLES_TEXT_OUTPUT string = `CSV                        PACKAGE         CHANNEL  REPLACES                   
prometheusoperator.0.47.0  prometheus      beta     prometheusoperator.0.37.0  
redis-operator.v0.8.0      redis-operator  stable   redis-operator.v0.6.0      
`

func TestBundlesCmdToText(t *testing.T) {
	out, err := bundlesCmdToText(EXPECTED_BUNDLES)
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_BUNDLES_TEXT_OUTPUT, out)
}

const EXPECTED_BUNDLES_JSON_OUTPUT string = `[{"csvName":"prometheusoperator.0.47.0","packageName":"prometheus","channelname":"beta","bundlePath":"quay.io/openshift-community-operators/prometheus:v0.47.0","version":"0.47.0","replaces":"prometheusoperator.0.37.0"},{"csvName":"redis-operator.v0.8.0","packageName":"redis-operator","channelname":"stable","bundlePath":"quay.io/openshift-community-operators/redis-operator:v0.8.0","version":"0.8.0","replaces":"redis-operator.v0.6.0"}]`

func TestBundlesCmdToJson(t *testing.T) {
	out, err := bundlesCmdToJson(EXPECTED_BUNDLES)
	assert.Nil(t, err)
	assert.Equal(t, EXPECTED_BUNDLES_JSON_OUTPUT, out)
}
