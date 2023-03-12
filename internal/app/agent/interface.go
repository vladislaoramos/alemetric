package agent

import "github.com/vladislaoramos/alemetric/internal/entity"

type WebAPIAgent interface {
	SendMetrics(string, string, *entity.Counter, *entity.Gauge) error
	SendSeveralMetrics([]entity.Metrics) error
}
