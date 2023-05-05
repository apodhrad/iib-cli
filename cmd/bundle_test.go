package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBundle(t *testing.T) {
	stdout, stderr, err := testCmd(t, "get", "bundle", "prometheus", "beta", "prometheusoperator.0.56.3")
	assert.Nil(t, err)
	expected := readTestResource(t, "bundle.txt")
	assert.Equal(t, expected, stdout)
	assert.Equal(t, "", stderr)
}

func TestGetBundleJson(t *testing.T) {
	stdout, stderr, err := testCmd(t, "get", "bundle", "prometheus", "beta", "prometheusoperator.0.56.3", "-o", "json")
	assert.Nil(t, err)
	expected := readTestResource(t, "bundle.json")
	assert.Equal(t, expected, stdout)
	assert.Equal(t, "", stderr)
}
