package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	stdout, stderr, err := testCmd(t, "get")
	assert.Nil(t, err)
	expected := readTestResource(t, "get.stdout")
	assert.Equal(t, expected, stdout)
	expected = readTestResource(t, "get.stderr")
	assert.Equal(t, expected, stderr)
}

func TestGetHelp(t *testing.T) {
	stdout, stderr, err := testCmd(t, "get", "-h")
	assert.Nil(t, err)
	expected := readTestResource(t, "get-help.stdout")
	assert.Equal(t, expected, stdout)
	expected = readTestResource(t, "get-help.stderr")
	assert.Equal(t, expected, stderr)
}
