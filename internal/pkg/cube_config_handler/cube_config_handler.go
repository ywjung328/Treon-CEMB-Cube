package cube_config_handler

import (
	"os"

	"gopkg.in/yaml.v2"
)

func ReadConfig(path string) (CubeConfig, error) {
	yamlFile, err := os.ReadFile(path)
	var config CubeConfig
	err = yaml.Unmarshal(yamlFile, &config)

	if err != nil {
		return CubeConfig{}, err
	}

	return config, nil
}
