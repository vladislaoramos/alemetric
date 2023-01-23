package agent

import (
	"github.com/go-resty/resty/v2"
	"github.com/vladislaoramos/alemetric/configs"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Run(cfg *configs.Config) {
	lgr := logger.New(cfg.Logger.Level)

	metrics := NewMetrics()

	client := resty.New().SetBaseURL("http://" + net.JoinHostPort(cfg.Agent.Host, cfg.Agent.Port))

	webAPI := NewWebAPI(client)

	worker := NewWorker(lgr, metrics, cfg.Agent.MetricsNames, webAPI)

	updateTicker := time.NewTicker(time.Second * time.Duration(cfg.Agent.PollInterval))
	go worker.UpdateMetrics(updateTicker)

	sendTicker := time.NewTicker(time.Second * time.Duration(cfg.Agent.ReportInterval))
	go worker.SendMetrics(sendTicker)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	stop := <-sigs

	lgr.Info("agent got stop signal: " + stop.String())
}
