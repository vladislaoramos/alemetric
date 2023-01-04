package agent

import (
	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/pkg/log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *configs.Config) {
	l := logger.New(cfg.Logger.Level)

	metrics := NewMetrics()

	webAPI := NewAPI(l, cfg.Agent.Host, cfg.Agent.Port)

	worker := NewWorker(l, metrics, webAPI)

	updateTicker := time.NewTicker(time.Duration(cfg.Agent.PollInterval) * time.Second)
	go worker.UpdateMetrics(updateTicker)

	sendTicker := time.NewTicker(time.Duration(cfg.Agent.ReportInterval) * time.Second)
	go worker.SendMetrics(sendTicker)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	stop := <-sigs

	l.Info("agent got stop signal: " + stop.String())
}
