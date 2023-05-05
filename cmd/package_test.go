package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCmdGetPackage(t *testing.T) {
	stdout, stderr, err := testCmd(t, "get", "package", "redis-operator")
	assert.Nil(t, err)
	expected := readTestResource(t, "package.txt")
	assert.Equal(t, expected, stdout)
	assert.Equal(t, "", stderr)
}

func TestCmdGetPackageJson(t *testing.T) {
	stdout, stderr, err := testCmd(t, "get", "package", "redis-operator", "-o", "json")
	assert.Nil(t, err)
	expected := readTestResource(t, "package.json")
	assert.Equal(t, expected, stdout)
	assert.Equal(t, "", stderr)
}

func TestCmdGetPackageNone(t *testing.T) {
	stdout, stderr, err := testCmd(t, "get", "package")
	assert.NotNil(t, err)
	assert.Equal(t, "", stdout)
	expected := readTestResource(t, "package-none.stderr")
	assert.Equal(t, expected, stderr)
}

func TestCmdGetPackageIncorrect(t *testing.T) {
	stdout, stderr, err := testCmd(t, "get", "package", "foo-operator")
	assert.NotNil(t, err)
	assert.Equal(t, "", stdout)
	expected := readTestResource(t, "package-incorrect.stderr")
	assert.Equal(t, expected, stderr)
}
