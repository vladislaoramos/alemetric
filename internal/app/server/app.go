package server

import (
	"github.com/go-chi/chi/v5"
	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/repo"
	"github.com/vladislaoramos/alemetric/internal/usecase"
	"github.com/vladislaoramos/alemetric/pkg/log"
	"net/http"
)

func Run(cfg *configs.Config) {
	lgr := logger.New(cfg.Logger.Level)

	repoOpts := make([]repo.OptionFunc, 0)
	if cfg.Server.StoreFile != " " {
		repoOpts = append(repoOpts, repo.StoreFilePath(cfg.Server.StoreFile))
	}
	if cfg.Server.Restore && cfg.Server.StoreFile != " " {
		repoOpts = append(repoOpts, repo.Restore())
	}
	metricsRepo := repo.NewMetricsRepo(repoOpts...)

	mtOptions := make([]usecase.OptionFunc, 0)
	if cfg.Server.StoreInterval != 0 {
		mtOptions = append(mtOptions, usecase.WriteFileWithDuration(cfg.Server.StoreInterval))
	} else {
		mtOptions = append(mtOptions, usecase.SyncWriteFile())
	}

	handler := chi.NewRouter()

	mt := usecase.NewMetricsTool(metricsRepo, lgr, mtOptions...)
	NewRouter(handler, mt, lgr)

	lgr.Fatal(http.ListenAndServe(cfg.Address, handler).Error())
}
