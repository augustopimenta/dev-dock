package configs

import (
	"io/ioutil"
	"gopkg.in/yaml.v2"
	"os"
)

const configFile = "config.yaml"

type ConfigFile struct {
	Projects []Project `yaml:"projects"`
}

type Project struct {
	Name string `yaml:"name"`
	Domain string `yaml:"domain"`
	Image string `yaml:"image"`
	Volumes []string `yaml:"volumes"`
	Ports []string `yaml:"ports"`
	Envs []string `yaml:"envs"`
}

func (project Project) ToSlice() []string {
	return []string { project.Name, project.Domain, project.Image, "DOWN", "", "" }
}

func Exists() bool {
	if _, err := os.Stat(configFile); err == nil {
		return true
	}
	return false
}

func Create() {
	config := ConfigFile{[]Project{{}}}
	configYaml, _ := yaml.Marshal(config);
	err := ioutil.WriteFile(configFile, []byte(configYaml), 0644)
	if (err != nil) {
		panic(err)
	}
}

func Read() ConfigFile {
	data, err := ioutil.ReadFile(configFile)
	if (err != nil) {
		panic(err)
	}

	conf := ConfigFile{}
	err = yaml.Unmarshal(data, &conf)
	if (err != nil) {
		panic(err)
	}
	return conf
}