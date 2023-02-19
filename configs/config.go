package configs

import (
	"flag"
	"fmt"
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
	configPath = "./configs/config.yml"

	serverURL = "127.0.0.1:8080"

	pollInterval   = time.Second * 2
	reportInterval = time.Second * 10
	storeInterval  = time.Second * 300
	storeFile      = "/tmp/devops-metrics-db.json"

	agentName  = "alemetric-agent"
	serverName = "alemetric-server"

	AgentConfig  = "agent"
	ServerConfig = "server"

	loggerDefaultLevel = "debug"
)

var metricsNames = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
	"RandomValue",
	"PollCount",
}

func NewConfig(app string) *Config {
	envCfg := new(Config)
	_ = cleanenv.ReadEnv(envCfg)

	flagCfg := new(Config)

	switch app {
	case AgentConfig:
		flag.StringVar(&flagCfg.Agent.ServerURL, "a", serverURL, "server address")
		flag.DurationVar(&flagCfg.Agent.ReportInterval, "r", reportInterval, "report interval")
		flag.DurationVar(&flagCfg.Agent.PollInterval, "p", pollInterval, "poll interval")
		flag.StringVar(&flagCfg.Agent.Key, "k", flagCfg.Agent.Key, "encryption key")
	case ServerConfig:
		flag.StringVar(&flagCfg.Server.Address, "a", serverURL, "server address")
		flag.BoolVar(&flagCfg.Server.Restore, "r", true, "restore data from file")
		flag.DurationVar(&flagCfg.Server.StoreInterval, "i", storeInterval, "store interval")
		flag.StringVar(&flagCfg.Server.StoreFile, "f", storeFile, "store file")
		flag.StringVar(&flagCfg.Server.Key, "k", flagCfg.Server.Key, "encryption key")
		flag.StringVar(&flagCfg.Database.URL, "d", flagCfg.Database.URL, "database")
	}

	flag.Parse()

	if envCfg.Agent.ServerURL == "" {
		envCfg.Agent.ServerURL = flagCfg.Agent.ServerURL
	}

	if envCfg.Agent.ReportInterval.String() == "0s" {
		envCfg.Agent.ReportInterval = flagCfg.Agent.ReportInterval
	}

	if envCfg.Agent.PollInterval.String() == "0s" {
		envCfg.Agent.PollInterval = flagCfg.Agent.PollInterval
	}

	if envCfg.Agent.Key == "" {
		envCfg.Agent.Key = flagCfg.Agent.Key
	}

	if envCfg.Server.Address == "" {
		envCfg.Server.Address = flagCfg.Server.Address
	}

	if !envCfg.Server.Restore {
		envCfg.Server.Restore = flagCfg.Server.Restore
	}

	if envCfg.Server.StoreInterval.String() == "0s" {
		envCfg.Server.StoreInterval = flagCfg.Server.StoreInterval
	}

	if envCfg.Server.StoreFile == "" {
		envCfg.Server.StoreFile = flagCfg.Server.StoreFile
	}

	if envCfg.Server.Key == "" {
		envCfg.Server.Key = flagCfg.Server.Key
	}

	if envCfg.Database.URL == "" {
		envCfg.Database.URL = flagCfg.Database.URL
	}

	return envCfg
}

func (c Config) String() string {
	return fmt.Sprintf("restore: %v storeFile: %v storeInterval: %v", c.Restore, c.StoreFile, c.StoreInterval)
}