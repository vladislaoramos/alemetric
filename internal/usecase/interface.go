package usecase

import (
	"context"

	"github.com/vladislaoramos/alemetric/internal/entity"
)

// MetricsTool defines the interface of interaction between a client and the tool.
type MetricsTool interface {
	GetMetricsNames(context.Context) ([]string, error)
	StoreMetrics(context.Context, entity.Metrics) error
	GetMetrics(context.Context, entity.Metrics) (entity.Metrics, error)
	PingRepo(context.Context) error
}

// MetricsRepo defines the interface of interaction between the tool and the repository storage.
type MetricsRepo interface {
	StoreMetrics(context.Context, entity.Metrics) error
	GetMetrics(context.Context, string) (entity.Metrics, error)
	GetMetricsNames(ctx context.Context) []string
	StoreAll() error
	Upload(context.Context) error
	Ping(context.Context) error
}
