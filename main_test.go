package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMainFunc(t *testing.T) {
	out := testMain(t)
	expected := readTestResource(t, "main.txt")
	assert.Equal(t, expected, out)
}

func testMain(t *testing.T, cmdArgs ...string) string {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// this gets captured
	originalArgs := os.Args
	os.Args = []string{"iib-cli"}
	os.Args = append(os.Args, cmdArgs...)
	main()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	os.Args = originalArgs

	return string(out)
}

func readTestResource(t *testing.T, filename string) string {
	data, err := os.ReadFile(filename)
	assert.Nil(t, err)
	return string(data)
}
