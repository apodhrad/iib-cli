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
func setTestIIB(t *testing.T) {
	TestLogger.Println("Set IIB for test " + t.Name())
	os.Setenv("IIB", "quay.io/apodhrad/iib-test:v0.0.1")
}

func startTestGrpc(t *testing.T) {
	setTestIIB(t)
	grpc.GrpcStart()
}

func stopTestGrpc(t *testing.T) {
	os.Unsetenv("IIB")
	grpc.GrpcStop()
}

func testCmd(t *testing.T, cmdArgs ...string) (string, error) {
	setTestIIB(t)
	defer stopTestGrpc(t)

	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	// always reset the output
	output = ""

	// this gets captured
	originalArgs := os.Args
	os.Args = []string{"iib-cli"}
	os.Args = append(os.Args, cmdArgs...)
	err := rootCmd.Execute()

	w.Close()
	out, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	os.Args = originalArgs

	// v, _ := rootCmd.Flags().GetString("output")
	// fmt.Println(">>> output = " + v)

	return string(out), err
}

func readTestResource(t *testing.T, filename string) string {
	data, err := os.ReadFile("test-resources/" + filename)
	assert.Nil(t, err)
	return string(data)
}
