package server

import (
	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/repo"
	"github.com/vladislaoramos/alemetric/pkg/log"
	"net"
	"net/http"
)

func Run(cfg *configs.Config) {
	lgr := logger.New(cfg.Logger.Level)
	metricsRepo := repo.NewMetricsRepo()

	http.HandleFunc(
		"/update/",
		UpdateMetricsHandler(
			metricsRepo,
		),
	)

	lgr.Fatal(
		http.ListenAndServe(
			net.JoinHostPort(cfg.Server.Host, cfg.Server.Port), nil).
			Error())
}
