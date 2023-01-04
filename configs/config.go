package configs

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Agent  `yaml:"agent"`
	Logger `yaml:"logger"`
}

type Logger struct {
	Level string `yaml:"level"`
}

type Agent struct {
	Name           string `yaml:"name"`
	Version        string `yaml:"version"`
	PollInterval   int64  `yaml:"pollInterval"`
	ReportInterval int64  `yaml:"reportInterval"`
	Host           string `yaml:"host"`
	Port           string `yaml:"port"`
}

const configPath = "./configs/config.yml"

func NewConfig() (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, errors.New(err.Error())
	}

	return cfg, nil
}
