package benchmark

import (
	"context"
	"fmt"
	"testing"

	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/internal/repo"
)

func BenchmarkGetMetricsNames(b *testing.B) {
	metricsRepo, _ := repo.NewMetricsRepo()
	ctx := context.Background()

	for i := 0; i < 10000; i++ {
		metrics := entity.Metrics{
			ID:    fmt.Sprintf("metric%d", i),
			MType: "type",
		}
		_ = metricsRepo.StoreMetrics(ctx, metrics)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = metricsRepo.GetMetricsNames(context.Background())
	}
}

func BenchmarkStoreMetrics(b *testing.B) {
	metricsRepo, _ := repo.NewMetricsRepo()
	ctx := context.Background()

	metrics := entity.Metrics{
		ID:    "metric1",
		MType: "type",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = metricsRepo.StoreMetrics(ctx, metrics)
	}
}

func BenchmarkGetMetrics(b *testing.B) {
	metricsRepo, _ := repo.NewMetricsRepo()
	ctx := context.Background()

	for i := 0; i < 10000; i++ {
		metrics := entity.Metrics{
			ID:    fmt.Sprintf("metric%d", i),
			MType: "type",
		}
		_ = metricsRepo.StoreMetrics(ctx, metrics)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = metricsRepo.GetMetrics(ctx, "metric1")
	}
}

func BenchmarkStoreAll(b *testing.B) {
	metricsRepo, _ := repo.NewMetricsRepo()
	ctx := context.Background()

	// Заполняем хранилище данными для тестирования
	for i := 0; i < 10000; i++ {
		metrics := entity.Metrics{
			ID:    fmt.Sprintf("metric%d", i),
			MType: "type",
		}
		_ = metricsRepo.StoreMetrics(ctx, metrics)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = metricsRepo.StoreAll()
	}
}

func BenchmarkUpload(b *testing.B) {
	metricsRepo, _ := repo.NewMetricsRepo()
	ctx := context.Background()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = metricsRepo.Upload(ctx)
	}
}
