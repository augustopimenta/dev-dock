package configs

import (
	"os"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const (
	configFile = "config.yaml"

	ExampleProjectName = "example"
)

type ConfigFile struct {
	UseVirtualHost bool `yaml:"virtual-hosts"`
	Projects []Project `yaml:"projects"`
}

func (c ConfigFile) FindProject(name string) *Project {
	for _, project := range c.Projects {
		if project.Name == name && project.Name != ExampleProjectName {
			return &project
		}
	}

	return nil;
}

func NewConfigFile() ConfigFile {
	if exists() {
		return read()
	}

	return create()
}

func exists() bool {
	_, err := os.Stat(configFile)

	return err == nil
}

func read() ConfigFile {
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

func create() ConfigFile {

	exampleProject := Project{
		Name: ExampleProjectName,
		Domain: "optional.example.dev",
		Image: "ubuntu:16.04",
		Volumes: []string{"/users/user/app:/var/www/app"},
		Ports: []string{"3000:80"},
		Envs: []string{"SOME_VAR=1"},
	}

	config := ConfigFile{ true, []Project{exampleProject} }
	configYaml, _ := yaml.Marshal(config);
	err := ioutil.WriteFile(configFile, []byte(configYaml), 0644)
	if (err != nil) {
		panic(err)
	}

	return config;
}

