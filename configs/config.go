package configs

import (
	"os"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

const configFile = "config.yaml"

type ConfigFile struct {
	UseVirtualHost bool `yaml:"virtual-hosts"`
	Projects []Project `yaml:"projects"`
}

func (c ConfigFile) FindProject(name string) *Project {
	for _, project := range c.Projects {
		if project.Name == name {
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
	config := ConfigFile{}
	configYaml, _ := yaml.Marshal(config);
	err := ioutil.WriteFile(configFile, []byte(configYaml), 0644)
	if (err != nil) {
		panic(err)
	}

	return config;
}

