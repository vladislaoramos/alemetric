package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vladislaoramos/alemetric/configs"
	"github.com/vladislaoramos/alemetric/internal/repo"
	"github.com/vladislaoramos/alemetric/internal/usecase"
	"github.com/vladislaoramos/alemetric/pkg/log"
	"github.com/vladislaoramos/alemetric/pkg/postgres"
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

	mtOptions := make([]usecase.OptionFunc, 0)
	if cfg.Server.StoreInterval != 0 {
		mtOptions = append(mtOptions, usecase.WriteFileWithDuration(cfg.Server.StoreInterval))
	} else {
		mtOptions = append(mtOptions, usecase.SyncWriteFile())
	}
	if cfg.Server.Key != "" {
		mtOptions = append(mtOptions, usecase.CheckDataSign(cfg.Server.Key))
	}

	var curRepo usecase.MetricsRepo
	if cfg.Database.URL != "" {
		err := runMigration(cfg.Database.URL, cfg.Database.MigrationDir)
		if err != nil {
			lgr.Fatal(fmt.Sprintf("Server - Run Migration - Error: %s", err.Error()))
		}

		db, err := postgres.New(cfg.Database.URL)
		if err != nil {
			lgr.Fatal(fmt.Sprintf("Server - PostgreSQL Init - Error: %s", err.Error()))
		}
		defer db.Close()

		curRepo = repo.NewPostgresRepo(db)
	} else {
		curRepo = repo.NewMetricsRepo(repoOpts...)
	}

	handler := chi.NewRouter()

	mt := usecase.NewMetricsTool(curRepo, lgr, mtOptions...)
	NewRouter(handler, mt, lgr)

	lgr.Fatal(http.ListenAndServe(cfg.Address, handler).Error())
}
