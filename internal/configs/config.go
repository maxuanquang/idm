package configs

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type ConfigFilePath string

type Config struct {
	Auth     Auth     `yaml:"auth"`
	Database Database `yaml:"database"`
	Log      Log      `yaml:"log"`
}

func NewConfig(configFilePath ConfigFilePath) (Config, error) {
	configBytes, err := os.ReadFile(string(configFilePath))
	if err != nil {
		return Config{}, fmt.Errorf("error reading configuration file: %w", err)
	}

	config := Config{}
	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshal configuration file: %w", err)
	}

	return config, nil
}
