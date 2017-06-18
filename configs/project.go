package configs

import "strings"

type Project struct {
	Name string `yaml:"name"`
	Domain string `yaml:"domain"`
	Image string `yaml:"image"`
	Status string `yaml:"-"`
	Volumes []string `yaml:"volumes"`
	Ports []string `yaml:"ports"`
	Envs []string `yaml:"envs"`
}

func (project Project) ToSlice() []string {
	return []string {
		project.Name,
		project.Domain,
		project.Image,
		project.Status,
		strings.Join(project.Volumes, "\n"),
		strings.Join(project.Ports, "\n"),
	}
}
