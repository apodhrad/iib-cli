package grpc

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

// Initialize a docker client
func newClient() (context.Context, *client.Client) {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		ErrorLogger.Panicln(err)
	}
	return ctx, cli
}

// Pull a docker image
func DockerPullImage(image string) error {
	InfoLogger.Printf("Pull image %s\n", image)
	ctx, cli := newClient()
	defer cli.Close()

	DebugLogger.Printf("docker pull %s\n", image)
	out, err := cli.ImagePull(ctx, image, types.ImagePullOptions{})
	defer out.Close()

	if err != nil {
		InfoLogger.Printf("Image %s was successfully pulled\n", image)
	}
	return err
}

// Start a docker container
func DockerStartContainer(name string, image string, portMapping string) (string, error) {
	InfoLogger.Printf("Start a container with name %s, image %s and port mapping %s\n", name, image, portMapping)

	ctx, cli := newClient()
	defer cli.Close()

	containerConfig := &container.Config{Image: image}
	containerHostConfig := &container.HostConfig{}

	if portMapping != "" {
		exposedPort := strings.Split(portMapping, ":")[0]
		mappedPort := strings.Split(portMapping, ":")[1]

		containerConfig = &container.Config{
			Image:        image,
			ExposedPorts: nat.PortSet{nat.Port(exposedPort): struct{}{}},
		}

		containerHostConfig = &container.HostConfig{
			PortBindings: nat.PortMap{
				nat.Port(exposedPort): []nat.PortBinding{
					{
						HostIP:   "0.0.0.0",
						HostPort: mappedPort,
					},
				},
			},
		}

	}

	DebugLogger.Printf("docker create --name %s -p %s %s\n", name, portMapping, image)
	resp, err := cli.ContainerCreate(ctx, containerConfig, containerHostConfig, nil, nil, name)
	if err != nil {
		ErrorLogger.Println(err)
		return "", err
	}
	DebugLogger.Println("Wait for state [created]")
	c, err := waitForState(name, "created")
	if err != nil {
		ErrorLogger.Println(err)
		return "", err
	}
	DebugLogger.Println(containerToString(c))

	DebugLogger.Printf("docker start -d %s\n", name)
	err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
	if err != nil {
		ErrorLogger.Println(err)
		return "", err
	}
	DebugLogger.Println("Wait for state [running]")
	c, err = waitForState(name, "running")
	if err != nil {
		ErrorLogger.Println(err)
		return "", err
	}
	InfoLogger.Println(containerToString(c))

	return resp.ID, nil
}

// Stop a docker container
func DockerStopContainer(name string) error {
	InfoLogger.Printf("Stop a container with name '%s'\n", name)

	// check if there already is a container
	c := getContainerWithIdOrName(name)
	if c == nil {
		InfoLogger.Printf("No container with name '%s' found\n", name)
		return nil
	}
	InfoLogger.Printf("Container with name '%s' found\n", name)

	// now, wait for a proper state before we stop the container
	DebugLogger.Println("Wait for a proper state [created|running|exited]")
	c, err := waitForState(name, "created", "running", "exited")
	if err != nil {
		ErrorLogger.Println(err)
		return err
	}

	ctx, cli := newClient()
	defer cli.Close()

	// stop the container and wait for the state [exited]
	DebugLogger.Printf("docker stop %s\n", name)
	err = cli.ContainerStop(ctx, c.ID, container.StopOptions{})
	if err != nil {
		ErrorLogger.Println(err)
		return err
	}
	DebugLogger.Println("Wait for state [exited]")
	c, err = waitForState(name, "exited")
	if err != nil {
		ErrorLogger.Println(err)
		return err
	}
	DebugLogger.Println(containerToString(c))

	// once the conatainer is in exited state, we can remove it
	DebugLogger.Printf("docker rm %s\n", name)
	err = cli.ContainerRemove(ctx, c.ID, types.ContainerRemoveOptions{})
	if err != nil {
		ErrorLogger.Println(err)
		return err
	}
	DebugLogger.Println("Wait until the container is gone")
	c, err = waitForState(name)
	if err != nil {
		ErrorLogger.Println(err)
		return err
	}

	InfoLogger.Printf("Container with name %s was successfully stopped and removed\n", name)
	return nil
}

func containerToString(c *types.Container) string {
	if c == nil {
		return "Container[nil]"
	}
	return fmt.Sprintf("Container[Name: %s, Image: %s, State: %s, Status: %s]", c.Names[0], c.Image, c.State, c.Status)
}

func getContainerWithIdOrName(idOrName string) *types.Container {
	ctx, cli := newClient()
	defer cli.Close()

	containers, err := cli.ContainerList(ctx, types.ContainerListOptions{All: true})
	if err != nil {
		ErrorLogger.Panicln(err)
	}

	for _, container := range containers {
		if container.ID == idOrName {
			return &container
		}
		for _, name := range container.Names {
			if name == "/"+idOrName {
				return &container
			}
		}
	}

	return nil
}

func waitForState(name string, states ...string) (*types.Container, error) {
	var c *types.Container

	for i := 0; i < 10; i++ {
		c = getContainerWithIdOrName(name)
		if len(states) > 0 {
			for _, state := range states {
				if c != nil && c.State == state {
					return c, nil
				}
			}
		} else {
			if c == nil {
				return c, nil
			}
		}
		time.Sleep(1 * time.Second)
	}

	if len(states) > 0 {
		return nil, fmt.Errorf("%s is still not in a state %s", containerToString(c), states)
	} else {
		return nil, fmt.Errorf("%s is still present", containerToString(c))
	}
}
