package configs

import (
	"flag"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
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
	Name           string        `yaml:"name"`
	PollInterval   time.Duration `yaml:"pollInterval" env:"POLL_INTERVAL"`
	ReportInterval time.Duration `yaml:"reportInterval" env:"REPORT_INTERVAL"`
	ServerURL      string        `yaml:"serverURL" env:"ADDRESS"`
	MetricsNames   []string      `yaml:"metricsNames"`
}

type Server struct {
	Name          string        `yaml:"name"`
	Address       string        `yaml:"address" env:"ADDRESS"`
	StoreInterval time.Duration `yaml:"storeInterval" env:"STORE_INTERVAL"`
	StoreFile     string        `yaml:"storeFile" env:"STORE_FILE"`
	Restore       bool          `yaml:"restore" env:"RESTORE"`
}

const configPath = "./configs/config.yml"

func NewConfig() (*Config, error) {
	cfg := new(Config)

	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %s", err.Error())
	}

	if err = cleanenv.ReadEnv(cfg); err != nil {
		return nil, fmt.Errorf("error setting envs: %w", err)
	}

	return cfg, nil
}

func Init(cfg *Config) error {
	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return fmt.Errorf("error reading config: %s", err.Error())
	}

	flag.Parse()

	if err = cleanenv.ReadEnv(cfg); err != nil {
		return fmt.Errorf("error setting envs: %w", err)
	}

	return nil
}
