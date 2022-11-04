package docker

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

var Client = &DockerClient{}

type DockerClient struct {
	client *client.Client
}

func (c *DockerClient) preCall() error {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return err
	}

	c.client = client

	return nil
}

func (c *DockerClient) FindContainers() ([]string, error) {
	err := c.preCall()
	if err != nil {
		return nil, err
	}

	containers, err := c.client.ContainerList(context.Background(), types.ContainerListOptions{
		All: true,
	})
	if err != nil {
		return nil, err
	}

	var containerIDs []string
	for _, container := range containers {
		if !strings.Contains(container.Names[0], "thing__") {
			continue
		}

		containerIDs = append(containerIDs, container.ID)
	}

	return containerIDs, nil
}

func (c *DockerClient) StopContainer(containerID string) error {
	err := c.preCall()
	if err != nil {
		return err
	}

	err = c.client.ContainerStop(context.Background(), containerID, nil)
	if err != nil {
		return err
	}

	fmt.Println("Stopped container:", containerID)
	return nil
}

func (c *DockerClient) StopContainers(containerIDs []string) error {
	for _, containerID := range containerIDs {
		err := c.StopContainer(containerID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *DockerClient) RemoveContainer(containerID string) error {
	err := c.preCall()
	if err != nil {
		return err
	}

	err = c.client.ContainerRemove(context.Background(), containerID, types.ContainerRemoveOptions{})
	if err != nil {
		return err
	}

	fmt.Println("Removed container:", containerID)
	return nil
}

func (c *DockerClient) RemoveContainers(containerIDs []string) error {
	for _, containerID := range containerIDs {
		err := c.RemoveContainer(containerID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *DockerClient) RemoveAllContainers() error {
	containerIDs, err := c.FindContainers()
	if err != nil {
		return err
	}

	for _, containerID := range containerIDs {
		err := c.StopContainer(containerID)
		if err != nil {
			fmt.Println(err)
		}
	}

	return c.RemoveContainers(containerIDs)
}

func (c *DockerClient) Start(image, imageName, name string, env []string, portMap nat.PortMap) (string, error) {
	err := c.preCall()
	if err != nil {
		return "", err
	}

	if portMap == nil {
		portMap = nat.PortMap{}
	}

	r, err := c.client.ImagePull(context.Background(), image, types.ImagePullOptions{})
	if err != nil {
		return "", err
	}
	defer r.Close()
	_, err = io.ReadAll(r)
	if err != nil {
		return "", err
	}

	container, err := c.client.ContainerCreate(context.Background(), &container.Config{
		Image: imageName,
		Env:   env,
	}, &container.HostConfig{
		PortBindings: portMap,
	}, nil, nil, fmt.Sprintf("thing__%s", name))
	if err != nil {
		return "", err
	}

	err = c.client.ContainerStart(context.Background(), container.ID, types.ContainerStartOptions{})

	return container.ID, err
}
