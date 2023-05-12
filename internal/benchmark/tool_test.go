package benchmark

import (
	"bytes"
	"context"
	"testing"

	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/internal/repo"
	"github.com/vladislaoramos/alemetric/internal/usecase"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

func BenchmarkToolGetMetrics(b *testing.B) {
	metricsRepo, _ := repo.NewMetricsRepo()
	ctx := context.Background()
	var buf bytes.Buffer
	log := logger.New("debug", &buf)
	useCase := usecase.NewMetricsTool(metricsRepo, log)

	metrics := entity.Metrics{
		ID:    "metric1",
		MType: "type",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = useCase.GetMetrics(ctx, metrics)
	}
}

func BenchmarkToolGetMetricsNames(b *testing.B) {
	metricsRepo, _ := repo.NewMetricsRepo()
	ctx := context.Background()
	var buf bytes.Buffer
	log := logger.New("debug", &buf)
	useCase := usecase.NewMetricsTool(metricsRepo, log)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = useCase.GetMetricsNames(ctx)
	}
}

func BenchmarkToolStoreMetrics(b *testing.B) {
	metricsRepo, _ := repo.NewMetricsRepo()
	ctx := context.Background()
	var buf bytes.Buffer
	log := logger.New("debug", &buf)
	useCase := usecase.NewMetricsTool(metricsRepo, log)

	metrics := entity.Metrics{
		ID:    "metric1",
		MType: "type",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = useCase.StoreMetrics(ctx, metrics)
	}
}
