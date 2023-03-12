package usecase

import (
	"context"

	"github.com/vladislaoramos/alemetric/internal/entity"
)

type MetricsTool interface {
	GetMetricsNames(context.Context) ([]string, error)
	StoreMetrics(context.Context, entity.Metrics) error
	GetMetrics(context.Context, entity.Metrics) (entity.Metrics, error)
	PingRepo(context.Context) error
}

type MetricsRepo interface {
	StoreMetrics(context.Context, entity.Metrics) error
	GetMetrics(context.Context, string) (entity.Metrics, error)
	GetMetricsNames(ctx context.Context) []string
	StoreAll() error
	Upload(context.Context) error
	Ping(context.Context) error
}
