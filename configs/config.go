package configs

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Agent  `yaml:"agent"`
	Server `yaml:"server"`
	Logger `yaml:"logger"`
}

type Logger struct {
	Level string `yaml:"level"`
}

type Agent struct {
	Name           string   `yaml:"name"`
	PollInterval   int64    `yaml:"pollInterval"`
	ReportInterval int64    `yaml:"reportInterval"`
	Host           string   `yaml:"host"`
	Port           string   `yaml:"port"`
	MetricsNames   []string `yaml:"metricsNames"`
}

type Server struct {
	Name string `yaml:"name"`
	Host string `yaml:"host"`
	Port string `yaml:"port"`
}

const configPath = "./configs/config.yml"

func NewConfig() (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %s", err.Error())
	}

	return cfg, nil
}
