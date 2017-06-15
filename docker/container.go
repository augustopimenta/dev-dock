package docker

type Container struct {
	Id string
	Name string
	Image string
	Ports map[string]string
	Volumes []string
	Envs []string
}