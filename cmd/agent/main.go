package main

import (
	"fmt"
	"os"

	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/app/agent"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

func main() {
	aCfg := configs.NewConfig(configs.AgentConfig)
	lgr := logger.New(aCfg.Logger.Level, os.Stdout)
	lgr.Info(fmt.Sprintf("%+v", *aCfg))

	agent.Run(aCfg, lgr)
}
