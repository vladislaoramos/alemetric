package main

import (
	"fmt"
	"os"

	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/app/agent"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

func main() {
	agentCfg := configs.NewConfig(configs.AgentConfig)
	lgr := logger.New(agentCfg.Logger.Level, os.Stdout)
	lgr.Info(fmt.Sprintf("%+v", *agentCfg))

	agent.Run(agentCfg, lgr)
}
