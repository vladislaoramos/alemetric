package main

import (
	"flag"
	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/app/server"
	"log"
)

func main() {
	cfg := new(configs.Config)

	flag.StringVar(&cfg.Server.Address, "a", cfg.Server.Address, "server address")
	flag.BoolVar(&cfg.Server.Restore, "r", cfg.Server.Restore, "restore data from file")
	flag.DurationVar(&cfg.Server.StoreInterval, "i", cfg.Server.StoreInterval, "store interval")
	flag.StringVar(&cfg.Server.StoreFile, "f", cfg.Server.StoreFile, "store file")
	flag.StringVar(&cfg.Server.Key, "k", cfg.Server.Key, "encryption key")
	flag.StringVar(&cfg.Database.URL, "d", cfg.Database.URL, "database")

	err := configs.Init(cfg)
	if err != nil {
		log.Fatalf("Server - Config Init - Error: %s", err.Error())
	}

	server.Run(cfg)
}
