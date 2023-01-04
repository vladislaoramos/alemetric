package server

import (
	"github.com/vladislaoramos/alemetric/internal/entity"
)

type MetricsRepo interface {
	StoreGaugeMetrics(string, entity.Gauge) error
	StoreCounterMetrics(string, entity.Counter) error
	GetMetricsNames() []string
	GetMetrics(string) (interface{}, error)
}
