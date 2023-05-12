package main

import (
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"

	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/app/server"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

func main() {
	serverCfg := configs.NewConfig(configs.ServerConfig)

	f, err := os.OpenFile("/tmp/log_server", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	if err != nil {
		log.Fatal("unable to open file for log")
	}

	lgr := logger.New(serverCfg.Logger.Level, f)
	lgr.Info(fmt.Sprintf("%+v", serverCfg))

	server.Run(serverCfg, lgr)
}
