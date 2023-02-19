package main

import (
	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/app/server"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

func main() {
	cfg := configs.NewConfig(configs.ServerConfig)
	lgr := logger.New(cfg.Logger.Level)
	lgr.Info(cfg.String())
	server.Run(cfg)
}
