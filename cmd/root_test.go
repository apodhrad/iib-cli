package cmd

import (
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/apodhrad/iib-cli/grpc"
	"github.com/stretchr/testify/assert"
)

var (
	TestLogger *log.Logger
)

func init() {
	logFile := "/tmp/logs.txt"
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	TestLogger = log.New(file, "TEST: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func setup(t *testing.T) {
	TestLogger.Println("Set IIB for test " + t.Name())
	os.Setenv("IIB", "quay.io/apodhrad/iib-test:v0.0.1")
	grpc.GrpcStart()
}

func teardown(t *testing.T) {
	os.Unsetenv("IIB")
	grpc.GrpcStop()
}

func testCmd(t *testing.T, cmdArgs ...string) (string, string, error) {
	setup(t)
	defer teardown(t)

	// always reset the output
	output = ""

	// set args
	originalArgs := os.Args
	os.Args = []string{"iib-cli"}
	os.Args = append(os.Args, cmdArgs...)

	// catch stdout
	rescueStdout := os.Stdout
	rescueStderr := os.Stderr
	rStdout, wStdout, _ := os.Pipe()
	rStderr, wStderr, _ := os.Pipe()
	os.Stdout = wStdout
	os.Stderr = wStderr

	// this will be captured
	err := rootCmd.Execute()

	wStdout.Close()
	wStderr.Close()
	stdout, _ := ioutil.ReadAll(rStdout)
	stderr, _ := ioutil.ReadAll(rStderr)
	os.Stdout = rescueStdout
	os.Stderr = rescueStderr

	os.Args = originalArgs

	return string(stdout), string(stderr), err
}

func readTestResource(t *testing.T, filename string) string {
	data, err := os.ReadFile("test-resources/" + filename)
	assert.Nil(t, err)
	return string(data)
}
