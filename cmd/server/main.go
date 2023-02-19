package main

import (
	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/app/server"
	"log"
)

func main() {
	cfg, err := configs.NewConfig(configs.ServerConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	server.Run(cfg)
}
