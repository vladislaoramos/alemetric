package agent

import (
	"context"
	"github.com/vladislaoramos/alemetric/internal/entity"
)

// WebAPIAgent defines the interface of interaction between the agent and the server.
type WebAPIAgent interface {
	SendMetrics(context.Context, string, string, *entity.Counter, *entity.Gauge) error
	SendSeveralMetrics(context.Context, []entity.Metrics) error
	Connect() (func(), error)
}
