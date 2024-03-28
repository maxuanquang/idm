package configs

import (
	"fmt"
	"os"

	"github.com/maxuanquang/idm/configs"
	"gopkg.in/yaml.v2"
)

type ConfigFilePath string

type Config struct {
	Auth     Auth     `yaml:"auth"`
	Database Database `yaml:"database"`
	Log      Log      `yaml:"log"`
	Cache    Cache    `yaml:"cache"`
	HTTP     HTTP     `yaml:"http"`
	GRPC     GRPC     `yaml:"grpc"`
	MQ       MQ       `yaml:"mq"`
	Download Download `yaml:"download"`
}

func NewConfig(configFilePath ConfigFilePath) (Config, error) {
	var (
		configBytes []byte = configs.DefaultConfigBytes
		config      Config
		err         error
	)

	if configFilePath != "" {
		configBytes, err = os.ReadFile(string(configFilePath))
		if err != nil {
			return Config{}, fmt.Errorf("error reading configuration file: %w", err)
		}
	}

	err = yaml.Unmarshal(configBytes, &config)
	if err != nil {
		return Config{}, fmt.Errorf("error unmarshal configuration file: %w", err)
	}

	return config, nil
}
