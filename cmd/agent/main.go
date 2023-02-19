package main

import (
	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/app/agent"
)

func main() {
	agent.Run(configs.NewConfig(configs.AgentConfig))
}
