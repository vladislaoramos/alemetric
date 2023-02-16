package usecase

import (
	"context"
	"github.com/vladislaoramos/alemetric/internal/entity"
)

type MetricsTool interface {
	GetMetricsNames() ([]string, error)
	StoreMetrics(entity.Metrics) error
	GetMetrics(entity.Metrics) (entity.Metrics, error)
	PingRepo(context.Context) error
}

type MetricsRepo interface {
	StoreMetrics(entity.Metrics) error
	GetMetrics(string) (entity.Metrics, error)
	GetMetricsNames() []string
	StoreAll() error
	Upload() error
	Ping(context.Context) error
}
