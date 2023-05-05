package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdGetPackage(t *testing.T) {
	out, err := testCmd(t, "get", "package", "redis-operator")
	assert.Nil(t, err)
	expected := readTestResource(t, "package.txt")
	assert.Equal(t, expected, out)
}

func TestCmdGetPackageJson(t *testing.T) {
	out, err := testCmd(t, "get", "package", "redis-operator", "-o", "json")
	assert.Nil(t, err)
	expected := readTestResource(t, "package.json")
	assert.Equal(t, expected, out)
}

func TestCmdGetPackageNone(t *testing.T) {
	out, err := testCmd(t, "get", "package")

	assert.NotNil(t, err)
	assert.Equal(t, "", out)
}
