package grpc

import (
	"fmt"
	"os"
	"time"

	"github.com/apodhrad/iib-cli/logging"
)

const GRPC_NAME string = "iib_registry_server"
const GRPC_HOST string = "localhost"
const GRPC_PORT string = "50051"
const GRPC_SERVER string = GRPC_HOST + ":" + GRPC_PORT

const GRPC_STARTSTOP_ENV_KEY string = "GRPC_STARTSTOP"

func GrpcStart() string {
	logging.INFO().Printf("Start the server")

	if isStartStopDisabled() {
		logging.INFO().Printf("Start/Stop feature is disabled")
		err := waitForReadiness(GRPC_SERVER)
		if err == nil {
			logging.INFO().Printf("Sever is up and running on %s", GRPC_SERVER)
			return GRPC_SERVER
		}
		logging.INFO().Printf("Server doesn't respond, so let's start it as usual")
	}

	// check is iib was specified
	iib := os.Getenv("IIB")
	if iib == "" {
		handlePanic(fmt.Errorf("Specify index image via envvar IIB or via command set iib!"))
	}

	// make sure there is no running container before starting a new one
	logging.INFO().Printf("Make sure there is no running container")
	GrpcStop()

	// pull the index image bundle (iib)
	logging.INFO().Printf("Pull image %s\n", iib)
	err := DockerPullImage(iib)
	if err != nil {
		handlePanic(err)
	}

	// now, we can start the new container
	logging.INFO().Printf("Start grpc server on localhost:" + GRPC_PORT)
	id, err := DockerStartContainer(GRPC_NAME, iib, GRPC_PORT+":"+GRPC_PORT)
	if err != nil {
		handlePanic(err)
	}
	logging.INFO().Printf("Container wit ID %s was sucesfully started", id)

	// the container should be running, now wait for its readiness
	logging.INFO().Printf("Wait for its readiness")
	err = waitForReadiness(GRPC_SERVER)
	if err != nil {
		handlePanic(err)
	}
	logging.INFO().Printf("The grpc server is up and running on %s", GRPC_SERVER)
	return GRPC_SERVER
}

func GrpcStop() {
	logging.INFO().Printf("Stop the server")

	if isStartStopDisabled() {
		logging.INFO().Printf("Start/Stop feature is disabled")
		logging.INFO().Printf("Thus, the server will not be stopped")
		return
	}

	err := DockerStopContainer(GRPC_NAME)
	if err != nil {
		handlePanic(err)
	}
	logging.INFO().Printf("The grpc server is stopped")
}

func isStartStopDisabled() bool {
	startStop := os.Getenv(GRPC_STARTSTOP_ENV_KEY)
	return startStop == "false"
}

func waitForReadiness(address string) error {
	for i := 0; i < 10; i++ {
		client, err := NewClient(address)
		defer client.Close()
		if client != nil && err == nil {
			isReady, _ := client.HealthCheck()
			if isReady {
				return nil
			}
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("Server on %s is still not ready!", address)
}

func handlePanic(err error) {
	logging.ERROR().Println(err)
	panic(err)
}
