package main

import (
	"fmt"
	"github.com/vladislaoramos/alemetric/internal/app/agent"
	_ "net/http/pprof"
	"os"

	"github.com/vladislaoramos/alemetric/configs"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	printBuildInfo()

	agentCfg := configs.NewConfig(configs.AgentConfig)
	lgr := logger.New(agentCfg.Logger.Level, os.Stdout)
	lgr.Info(fmt.Sprintf("%+v", *agentCfg))

	agent.Run(agentCfg, lgr)
}

func printBuildInfo() {
	version := fmt.Sprintf("Build version: %s\n", buildVersion)
	data := fmt.Sprintf("Build date: %s\n", buildDate)
	commit := fmt.Sprintf("Build commit: %s\n", buildCommit)
	fmt.Print(version, data, commit)
}
