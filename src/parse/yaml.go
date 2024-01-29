package parse

import (
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type DevContainerConfig struct {
	ReposPath  []string            `yaml:"reposPath"`
	Containers map[string][]string `yaml:"containers"`
}

type YAMLConfig struct {
	DevContainers DevContainerConfig `yaml:"devContainers"`
}

var once sync.Once
var configInstance YAMLConfig

func parseYaml() (YAMLConfig, error) {
	yamlFile, err := os.ReadFile("configuration.yaml")

	if err != nil {
		return YAMLConfig{}, err
	}

	var yamlConfig YAMLConfig
	err = yaml.Unmarshal(yamlFile, &yamlConfig)

	if err != nil {
		return YAMLConfig{}, err
	}

	return yamlConfig, nil
}

func GetConfiguration() YAMLConfig {
	once.Do(func() {
		configInstance, _ = parseYaml()
	})

	return configInstance
}
