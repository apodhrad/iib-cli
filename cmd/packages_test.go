package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdGetPackages(t *testing.T) {
	stdout, stderr, err := testCmd(t, "get", "packages")
	assert.Nil(t, err)
	expected := readTestResource(t, "packages.txt")
	assert.Equal(t, expected, stdout)
	assert.Equal(t, "", stderr)
}

func TestCmdGetPackagesJson(t *testing.T) {
	stdout, stderr, err := testCmd(t, "get", "packages", "-o", "json")
	assert.Nil(t, err)
	expected := readTestResource(t, "packages.json")
	assert.Equal(t, expected, stdout)
	assert.Equal(t, "", stderr)
}
