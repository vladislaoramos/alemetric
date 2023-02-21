package main

import (
	"fmt"

	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/app/agent"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

func main() {
	srvCfg := configs.NewConfig(configs.ServerConfig)
	lgr := logger.New(srvCfg.Logger.Level)
	lgr.Info(fmt.Sprintf("%+v", *srvCfg))

	agent.Run(configs.NewConfig(configs.AgentConfig))
}
