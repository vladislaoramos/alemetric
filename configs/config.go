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
	configPath = "./configs/config.yml"

	serverURL = "127.0.0.1:8080"

	pollInterval   = time.Second * 2
	reportInterval = time.Second * 10
	storeInterval  = time.Second * 300
	storeFile      = "/tmp/devops-metrics-db.json"
	restoreFlag    = true

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

func defaultServerCfg() *Config {
	return &Config{
		Server: Server{
			Name:          serverName,
			Address:       serverURL,
			StoreInterval: storeInterval,
			StoreFile:     storeFile,
			Restore:       restoreFlag,
		},
		Logger: Logger{Level: loggerDefaultLevel},
	}
}

func defaultAgentCfg() *Config {
	return &Config{
		Agent: Agent{
			Name:           agentName,
			PollInterval:   pollInterval,
			ReportInterval: reportInterval,
			ServerURL:      serverURL,
			MetricsNames:   metricsNames,
		},
		Logger: Logger{Level: loggerDefaultLevel},
	}
}

func (c *Config) updateAgentConfigs(v *Config) {
	if v.Agent.Name != "" && c.Agent.Name != v.Agent.Name {
		c.Agent.Name = v.Agent.Name
	}

	if v.ServerURL != "" && c.ServerURL != v.ServerURL {
		c.ServerURL = v.ServerURL
	}

	if v.PollInterval.String() != "0s" && c.PollInterval != v.PollInterval {
		c.PollInterval = v.PollInterval
	}

	if v.ReportInterval.String() != "0s" && c.ReportInterval != v.ReportInterval {
		c.ReportInterval = v.ReportInterval
	}

	if v.Level != "" && c.Level != v.Level {
		c.Level = v.Level
	}
}

func (c *Config) updateServerConfigs(v *Config) {
	if v.Server.Name != "" && c.Server.Name != v.Server.Name {
		c.Server.Name = v.Server.Name
	}

	if v.Address != "" && c.Address != v.Address {
		c.Address = v.Address
	}

	if v.StoreFile != "" && c.StoreFile != v.StoreFile {
		c.StoreFile = v.StoreFile
	}

	if v.StoreInterval.String() != "0s" && c.StoreInterval != v.StoreInterval {
		c.StoreInterval = v.StoreInterval
	}

	if c.Restore != v.Restore {
		c.Restore = v.Restore
	}

	if v.Level != "" && c.Level != v.Level {
		c.Level = v.Level
	}

	if v.Database.URL != "" && c.Database.URL != v.Database.URL {
		c.Database.URL = v.Database.URL
	}
}

func (c *Config) parseFlags(app string) {
	switch app {
	case AgentConfig:
		flag.StringVar(&c.Agent.ServerURL, "a", serverURL, "server address")
		flag.DurationVar(&c.Agent.ReportInterval, "r", reportInterval, "report interval")
		flag.DurationVar(&c.Agent.PollInterval, "p", pollInterval, "poll interval")
		flag.StringVar(&c.Agent.Key, "k", "", "encryption key")
	case ServerConfig:
		flag.StringVar(&c.Server.Address, "a", "", "server address")
		flag.BoolVar(&c.Server.Restore, "r", true, "restore data from file")
		flag.DurationVar(&c.Server.StoreInterval, "i", 0, "store interval")
		flag.StringVar(&c.Server.StoreFile, "f", "", "store file")
		flag.StringVar(&c.Server.Key, "k", "", "encryption key")
		flag.StringVar(&c.Database.URL, "d", "", "database")
	}

	flag.Parse()
}

func NewConfig(app string) *Config {
	var (
		cfg   *Config
		envs  *Config
		flags *Config
	)

	switch app {
	case AgentConfig:
		cfg = defaultAgentCfg()

		flags = new(Config)
		flags.parseFlags(AgentConfig)
		cfg.updateAgentConfigs(flags)

		envs = new(Config)
		_ = cleanenv.ReadEnv(envs)
		cfg.updateAgentConfigs(envs)
	case ServerConfig:
		cfg = defaultServerCfg()

		flags = new(Config)
		flags.parseFlags(ServerConfig)
		cfg.updateServerConfigs(flags)

		envs = new(Config)
		_ = cleanenv.ReadEnv(envs)
		cfg.updateServerConfigs(envs)
	}

	return cfg
}
