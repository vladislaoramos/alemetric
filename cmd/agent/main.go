package main

import (
	"fmt"
	"os"

	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/app/agent"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

func main() {
	srvCfg := configs.NewConfig(configs.ServerConfig)
	lgr := logger.New(srvCfg.Logger.Level, os.Stdout)
	lgr.Info(fmt.Sprintf("%+v", *srvCfg))

	aCfg := configs.NewConfig(configs.AgentConfig)

	agent.Run(aCfg, lgr)
}
