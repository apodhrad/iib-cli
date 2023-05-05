package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetBundles(t *testing.T) {
	stdout, stderr, err := testCmd(t, "get", "bundles")
	assert.Nil(t, err)
	expected := readTestResource(t, "bundles.txt")
	assert.Equal(t, expected, stdout)
	assert.Equal(t, "", stderr)
}

func TestGetBundlesJson(t *testing.T) {
	stdout, stderr, err := testCmd(t, "get", "bundles", "-o", "json")
	assert.Nil(t, err)
	expected := readTestResource(t, "bundles.json")
	assert.Equal(t, expected, stdout)
	assert.Equal(t, "", stderr)
}
