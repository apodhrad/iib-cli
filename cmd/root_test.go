package cmd

import (
	"os"

	"github.com/apodhrad/iib-cli/utils"
)

func setup() {
	utils.GrpcStopSafely()
	os.Setenv("IIB", "quay.io/apodhrad/iib-test:v0.0.1")
	utils.GrpcStartSafely()
}

func teardown() {
	utils.GrpcStopSafely()
}
