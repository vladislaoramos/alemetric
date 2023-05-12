package agent

import "github.com/vladislaoramos/alemetric/internal/entity"

// WebAPIAgent defines the interface of interaction between the agent and the server.
type WebAPIAgent interface {
	SendMetrics(string, string, *entity.Counter, *entity.Gauge) error
	SendSeveralMetrics([]entity.Metrics) error
}
