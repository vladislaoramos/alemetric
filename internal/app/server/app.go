package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/repo"
	"github.com/vladislaoramos/alemetric/pkg/log"
	"net"
	"net/http"
)

func Run(cfg *configs.Config) {
	lgr := logger.New(cfg.Logger.Level)
	metricsRepo := repo.NewMetricsRepo()
	handler := chi.NewRouter()

	NewRouter(handler, metricsRepo)

	// lgr.Fatal(http.ListenAndServe(net.JoinHostPort("", cfg.Server.Port), handler).Error())
	lgr.Fatal(http.ListenAndServe(net.JoinHostPort(cfg.Server.Host, cfg.Server.Port), handler).Error())
}
