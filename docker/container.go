package docker

import (
	"strings"

	"devdock/configs"
	"github.com/docker/docker/api/types"
)

type Container struct {
	Id string
	Name string
	Image string
	Ports map[string]string
	Volumes []string
	Envs []string
}


func fromProject(project configs.Project) *Container {
	ports := make(map[string]string);

	for _, port := range project.Ports {
		portsConfig := strings.Split(port, ":")
		if len(portsConfig) == 1 {
			ports[portsConfig[0] + "/tcp"] = ""
		} else {
			ports[portsConfig[0] + "/tcp"] = portsConfig[1]
		}
	}

	container := &Container{
		Name: project.Name,
		Image: project.Image,
		Volumes: project.Volumes,
		Envs: project.Envs,
		Ports: ports,
	}

	return container
}

func StartProxyContainer() {
	api := Api{}
	api.Init()

	container := &Container{
		"",
		"reverse-proxy",
		"jwilder/nginx-proxy",
		map[string]string{"80/tcp":"80"},
		[]string{"/var/run/docker.sock:/tmp/docker.sock:ro"},
		[]string{},
	}

	if api.Has(container) {
		api.Remove(container)
	}
	api.Run(container)
}

func GetProjectContainer(project configs.Project) *types.ContainerJSON {
	api := Api{}
	api.Init()

	container := fromProject(project);

	return api.Get(container);
}

func StartProjectContainer(project configs.Project) {
	api := Api{}
	api.Init()

	container := fromProject(project);

	if api.Has(container) {
		api.Remove(container)
	}
	api.Run(container)
}

func FinishProjectContainer(project configs.Project) {
	api := Api{}
	api.Init()

	container := fromProject(project);

	if api.Has(container) {
		api.Remove(container)
	}
}