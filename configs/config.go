package configs

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Agent    `yaml:"agent"`
	Server   `yaml:"server"`
	Logger   `yaml:"logger"`
	Database `yaml:"database"`
}

type Logger struct {
	Level string `yaml:"level"`
}

type Database struct {
	URL string `env:"DATABASE_DSN"`
}

type Agent struct {
	Name           string        `yaml:"name"`
	PollInterval   time.Duration `yaml:"pollInterval" env:"POLL_INTERVAL"`
	ReportInterval time.Duration `yaml:"reportInterval" env:"REPORT_INTERVAL"`
	ServerURL      string        `yaml:"serverURL" env:"ADDRESS"`
	MetricsNames   []string      `yaml:"metricsNames"`
	Key            string        `env:"KEY"`
}

type Server struct {
	Name          string        `yaml:"name" env:"NAME"`
	Address       string        `yaml:"address" env:"ADDRESS"`
	StoreInterval time.Duration `yaml:"storeInterval" env:"STORE_INTERVAL"`
	StoreFile     string        `yaml:"storeFile" env:"STORE_FILE"`
	Restore       bool          `yaml:"restore" env:"RESTORE"`
	Key           string        `env:"KEY"`
}

const (
	configPath   = "./configs/config.yml"
	AgentConfig  = "agent"
	ServerConfig = "server"
)

func NewConfig(app string) (*Config, error) {
	cfg := new(Config)

	switch app {
	case AgentConfig:
		flag.StringVar(&cfg.Agent.ServerURL, "a", cfg.Agent.ServerURL, "server address")
		flag.DurationVar(&cfg.Agent.ReportInterval, "r", cfg.Agent.ReportInterval, "report interval")
		flag.DurationVar(&cfg.Agent.PollInterval, "p", cfg.Agent.PollInterval, "poll interval")
		flag.StringVar(&cfg.Agent.Key, "k", cfg.Agent.Key, "encryption key")
	case ServerConfig:
		flag.StringVar(&cfg.Server.Address, "a", cfg.Server.Address, "server address")
		flag.BoolVar(&cfg.Server.Restore, "r", cfg.Server.Restore, "restore data from file")
		flag.DurationVar(&cfg.Server.StoreInterval, "i", cfg.Server.StoreInterval, "store interval")
		flag.StringVar(&cfg.Server.StoreFile, "f", cfg.Server.StoreFile, "store file")
		flag.StringVar(&cfg.Server.Key, "k", cfg.Server.Key, "encryption key")
		flag.StringVar(&cfg.Database.URL, "d", cfg.Database.URL, "database")
	}

	// First: init from yaml
	err := cleanenv.ReadConfig(configPath, cfg)
	if err != nil {
		return nil, err
	}

	// Next: update from flags if there are any
	flag.Parse()

	// Finally: update from envs if there are any
	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
