package main

import (
	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/app/agent"
	"log"
)

func main() {
	cfg, err := configs.NewConfig(configs.AgentConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	agent.Run(cfg)
}
