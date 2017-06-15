package docker

import (
	"io"
	"os"
	"context"

	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
)

type Api struct {
	context context.Context
	client *client.Client
}

func (api *Api) Init() {
	api.context = context.Background()
	cli, err := client.NewEnvClient()
	if (err != nil) {
		panic(cli)
	}
	api.client = cli
}

func (api *Api) Run(c *Container) {
	api.pullImage(c.Image)
	api.createContainer(c)
	api.startContainer(c)
}

func (api *Api) Remove(c *Container) {

}

func (api *Api) pullImage(image string) {
	out, err := api.client.ImagePull(api.context, image, types.ImagePullOptions{})
	if (err != nil) {
		panic(err)
	}
	io.Copy(os.Stdout, out)
}

func (api *Api) findContainer(c *Container) {

}

func (api *Api) createContainer(c *Container) {

	exposedPorts := make(nat.PortSet)
	portBindings := make(nat.PortMap)

	for cPort, hPort := range c.Ports {
		exposedPorts[nat.Port(cPort)] = struct{}{}
		portBindings[nat.Port(cPort)] = []nat.PortBinding{{HostIP:"0.0.0.0", HostPort:hPort}}
	}

	containerConfigs := &container.Config{Image: c.Image, ExposedPorts: exposedPorts}

	hostConfigs := &container.HostConfig{Binds: c.Volumes, PortBindings: portBindings}

	resp, err := api.client.ContainerCreate(api.context, containerConfigs, hostConfigs, nil, "reverse-proxy")
	if err != nil {
		panic(err)
	}

	c.Id = resp.ID
}

func (api *Api) startContainer(c *Container) {
	if err := api.client.ContainerStart(api.context, c.Id, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}