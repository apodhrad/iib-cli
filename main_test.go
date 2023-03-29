package main

import (
	"os"
	"testing"
)

func TestMainApi(t *testing.T) {
	os.Args = append(os.Args, "api")
	os.Args = append(os.Args, "-h")
	main()
}
