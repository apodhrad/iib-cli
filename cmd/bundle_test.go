package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBundle(t *testing.T) {
	out, err := testCmd(t, "get", "bundle", "prometheus", "beta", "prometheusoperator.0.56.3")
	assert.Nil(t, err)
	expected := readTestResource(t, "bundle.txt")
	assert.Equal(t, expected, out)
}

func TestGetBundleJson(t *testing.T) {
	out, err := testCmd(t, "get", "bundle", "prometheus", "beta", "prometheusoperator.0.56.3", "-o", "json")
	assert.Nil(t, err)
	expected := readTestResource(t, "bundle.json")
	assert.Equal(t, expected, out)
}
