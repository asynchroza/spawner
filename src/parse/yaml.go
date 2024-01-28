package parse

import (
	"os"

	"gopkg.in/yaml.v2"
)

type YAMLConfig struct {
	DevContainers map[string][]string `yaml:"devContainers"`
}

func ParseYaml() (YAMLConfig, error) {
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
