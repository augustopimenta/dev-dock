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

func StartProxyContainer() {
	ctx := context.Background()
	console, err := client.NewEnvClient()
	if err != nil {
		panic(err)
	}

	out, err := console.ImagePull(ctx, "jwilder/nginx-proxy", types.ImagePullOptions{})
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, out);

	containerConfigs := &container.Config{
		Image: "jwilder/nginx-proxy",
		ExposedPorts: nat.PortSet{"80/tcp": {}},
	}

	hostConfigs := &container.HostConfig{
		Binds: []string{"/var/run/docker.sock:/tmp/docker.sock:ro"},
		PortBindings: nat.PortMap{"80/tcp":[]nat.PortBinding{
			{HostIP:"0.0.0.0", HostPort:"80"},
		}},
	}

	resp, err := console.ContainerCreate(ctx, containerConfigs, hostConfigs, nil, "reverse-proxy")
	if err != nil {
		panic(err)
	}

	if err := console.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}

	log, err := console.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		panic(err)
	}

	io.Copy(os.Stdout, log)
}
