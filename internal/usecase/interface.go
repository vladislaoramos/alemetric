package usecase

import (
	"github.com/vladislaoramos/alemetric/internal/entity"
)

type MetricsTool interface {
	GetMetricsNames() ([]string, error)
	StoreMetrics(entity.Metrics) error
	GetMetrics(entity.Metrics) (entity.Metrics, error)
}

type MetricsRepo interface {
	StoreMetrics(entity.Metrics) error
	GetMetrics(string) (entity.Metrics, error)
	GetMetricsNames() []string
	StoreToFile() error
	UploadFromFile() error
}
