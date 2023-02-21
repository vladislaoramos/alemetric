package main

import (
	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/app/agent"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

func main() {
	srvCfg := configs.NewConfig(configs.ServerConfig)
	lgr := logger.New(srvCfg.Logger.Level)
	lgr.Info(srvCfg.String())

	agent.Run(configs.NewConfig(configs.AgentConfig))
}
