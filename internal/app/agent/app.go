package agent

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/vladislaoramos/alemetric/configs"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

const urlProtocol = "http://"

func Run(cfg *configs.Config, lgr *logger.Logger) {
	metrics := NewMetrics()

	client := resty.New().SetBaseURL(urlProtocol + cfg.Agent.ServerURL)

	webAPI := NewWebAPI(client, cfg.Agent.Key)

	worker := NewWorker(lgr, metrics, cfg.Agent.MetricsNames, webAPI)

	updateTicker := time.NewTicker(cfg.Agent.PollInterval)
	go worker.UpdateMetrics(updateTicker)

	sendTicker := time.NewTicker(cfg.Agent.ReportInterval)
	go worker.SendMetrics(sendTicker)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	stop := <-sigs

	lgr.Info("agent got stop signal: " + stop.String())
}
