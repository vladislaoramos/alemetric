package agent

import (
	"github.com/go-resty/resty/v2"
	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/pkg/log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *configs.Config) {
	l := logger.New(cfg.Logger.Level)

	metrics := NewMetrics()

	client := resty.New().SetBaseURL("http://" + net.JoinHostPort(cfg.Agent.Host, cfg.Agent.Port))

	webAPI := NewAPI(client)

	worker := NewWorker(l, metrics, cfg.MetricsNames, webAPI)

	updateTicker := time.NewTicker(time.Duration(cfg.Agent.PollInterval) * time.Second)
	go worker.UpdateMetrics(updateTicker)

	sendTicker := time.NewTicker(time.Duration(cfg.Agent.ReportInterval) * time.Second)
	go worker.SendMetrics(sendTicker)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	stop := <-sigs

	l.Info("agent got stop signal: " + stop.String())
}
