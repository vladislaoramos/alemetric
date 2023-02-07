package main

import (
	"flag"
	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/app/agent"
	"log"
)

func main() {
	cfg := new(configs.Config)

	flag.StringVar(&cfg.Agent.ServerURL, "a", cfg.Agent.ServerURL, "server address")
	flag.DurationVar(&cfg.Agent.ReportInterval, "r", cfg.Agent.ReportInterval, "report interval")
	flag.DurationVar(&cfg.Agent.PollInterval, "p", cfg.Agent.PollInterval, "poll interval")
	flag.StringVar(&cfg.Agent.Key, "k", cfg.Agent.Key, "encryption key")

	err := configs.Init(cfg)
	if err != nil {
		log.Fatal(err.Error())
	}

	agent.Run(cfg)
}
