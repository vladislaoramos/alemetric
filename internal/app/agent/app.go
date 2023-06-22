// Package agent contains the implementation of the client application.
package agent

import (
	"github.com/vladislaoramos/alemetric/internal/app/agent/sender"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/vladislaoramos/alemetric/configs"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

const urlProtocol = "http://"

// Run method launches the client application.
func Run(cfg *configs.Config, lgr *logger.Logger) {
	metrics := NewMetrics()

	client := resty.New().SetBaseURL(urlProtocol + cfg.Agent.ServerURL)

	var webAPI WebAPIAgent
	if cfg.Agent.UseGRPC {
		webAPI = sender.NewGRPCAgent(cfg.ServerURL, client.BaseURL)
	} else {
		webAPI = sender.NewWebAPI(client, cfg.Agent.Key, cfg.Agent.CryptoKey)
	}

	worker := NewWorker(lgr, metrics, cfg.Agent.MetricsNames, webAPI, cfg.RateLimit)

	updateTicker := time.NewTicker(cfg.Agent.PollInterval)
	go worker.UpdateMetrics(updateTicker)
	go worker.UpdateAdditionalMetrics(updateTicker)

	sendTicker := time.NewTicker(cfg.Agent.ReportInterval)
	go worker.SendMetrics(sendTicker)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	stop := <-sigs

	lgr.Info("Agent got stop signal: " + stop.String())
}
