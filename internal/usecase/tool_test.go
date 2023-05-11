package usecase

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/internal/repo"
	"github.com/vladislaoramos/alemetric/internal/usecase/mocks"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
	"os"
	"testing"
)

func testLogger() *logger.Logger {
	f, err := os.OpenFile("/tmp/test_log_server", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	l := logger.New("debug", f)
	if err != nil {
		l.Fatal("unable to open file for log")
	}

	return l
}

func metricsTool(t *testing.T) (*ToolUseCase, *mocks.MetricsRepo) {
	log := testLogger()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockRepo := mocks.NewMetricsRepo(t)

	mockTool := NewMetricsTool(
		mockRepo,
		log,
	)

	return mockTool, mockRepo
}

func TestGetMetricsNames(t *testing.T) {
	tool, repoMock := metricsTool(t)
	ctx := context.Background()
	repoMock.On("GetMetricsNames", ctx).Return(nil)
	_, err := tool.GetMetricsNames(ctx)
	require.NoError(t, err)
}

func TestGetMetrics(t *testing.T) {
	t.Run("without encryption key", func(t *testing.T) {
		tool, repoMock := metricsTool(t)
		ctx := context.Background()
		metrics := entity.Metrics{ID: "id"}
		repoMock.On("GetMetrics", ctx, metrics.ID).Return(entity.Metrics{}, nil)
		_, err := tool.GetMetrics(ctx, metrics)
		require.NoError(t, err)
	})

	t.Run("with error not found", func(t *testing.T) {
		tool, repoMock := metricsTool(t)
		ctx := context.Background()
		metrics := entity.Metrics{ID: "id"}
		repoMock.On("GetMetrics", ctx, metrics.ID).Return(entity.Metrics{}, repo.ErrNotFound)
		_, err := tool.GetMetrics(ctx, metrics)
		require.Error(t, err)
	})

	t.Run("with error not found", func(t *testing.T) {
		tool, repoMock := metricsTool(t)
		ctx := context.Background()
		metrics := entity.Metrics{ID: "id"}
		repoMock.On("GetMetrics", ctx, metrics.ID).Return(entity.Metrics{}, errors.New("some error"))
		_, err := tool.GetMetrics(ctx, metrics)
		require.Error(t, err)
	})

	t.Run("with encryption key", func(t *testing.T) {
		tool, repoMock := metricsTool(t)
		ctx := context.Background()
		metrics := entity.Metrics{ID: "id"}
		repoMock.On("GetMetrics", ctx, metrics.ID).Return(entity.Metrics{}, nil)
		metrics, err := tool.GetMetrics(ctx, metrics)
		require.NoError(t, err)
		tool.encryptionKey = "key"
		require.NotNil(t, metrics.Hash)
	})
}

func TestPingRepo(t *testing.T) {
	t.Run("with error", func(t *testing.T) {
		tool, repoMock := metricsTool(t)
		ctx := context.Background()
		repoMock.On("Ping", ctx).Return(errors.New("some error"))
		err := tool.PingRepo(ctx)
		require.Error(t, err)
	})

	t.Run("without error", func(t *testing.T) {
		tool, repoMock := metricsTool(t)
		ctx := context.Background()
		repoMock.On("Ping", ctx).Return(nil)
		err := tool.PingRepo(ctx)
		require.NoError(t, err)
	})
}

func TestStoreMetrics(t *testing.T) {
	t.Run("error not implemented", func(t *testing.T) {
		tool, _ := metricsTool(t)
		ctx := context.Background()
		metricsGauge := entity.Metrics{ID: "id", MType: "some type"}
		err := tool.StoreMetrics(ctx, metricsGauge)
		require.Error(t, err)
		require.ErrorIs(t, err, ErrNotImplemented)
	})

	t.Run("gauge without error", func(t *testing.T) {
		tool, repoMock := metricsTool(t)
		ctx := context.Background()
		metricsGauge := entity.Metrics{ID: "id", MType: Gauge}
		repoMock.On("StoreMetrics", ctx, metricsGauge).Return(nil)
		err := tool.StoreMetrics(ctx, metricsGauge)
		require.NoError(t, err)
	})

	t.Run("gauge with some error", func(t *testing.T) {
		tool, repoMock := metricsTool(t)
		ctx := context.Background()
		metricsGauge := entity.Metrics{ID: "id", MType: Gauge}
		repoMock.On("StoreMetrics", ctx, metricsGauge).Return(errors.New("some error"))
		err := tool.StoreMetrics(ctx, metricsGauge)
		require.Error(t, err)
	})

	t.Run("gauge with error not found", func(t *testing.T) {
		tool, repoMock := metricsTool(t)
		ctx := context.Background()
		metricsGauge := entity.Metrics{ID: "id", MType: Gauge}
		repoMock.On("StoreMetrics", ctx, metricsGauge).Return(repo.ErrNotFound)
		err := tool.StoreMetrics(ctx, metricsGauge)
		require.Error(t, err)
	})

	t.Run("counter without error", func(t *testing.T) {
		tool, repoMock := metricsTool(t)
		ctx := context.Background()

		var delta entity.Counter = 5
		metricsCounter := entity.Metrics{ID: "id", MType: Counter, Delta: &delta}

		var old entity.Counter = 0
		repoMock.On("GetMetrics", ctx, metricsCounter.ID).Return(entity.Metrics{Delta: &old}, nil)
		repoMock.On("StoreMetrics", ctx, metricsCounter).Return(nil)

		err := tool.StoreMetrics(ctx, metricsCounter)
		require.NoError(t, err)
	})

	t.Run("counter with error not found", func(t *testing.T) {
		tool, repoMock := metricsTool(t)
		ctx := context.Background()

		var delta entity.Counter = 5
		metricsCounter := entity.Metrics{ID: "id", MType: Counter, Delta: &delta}

		var old entity.Counter = 0
		repoMock.On("GetMetrics", ctx, metricsCounter.ID).Return(entity.Metrics{Delta: &old}, repo.ErrNotFound)
		repoMock.On("StoreMetrics", ctx, metricsCounter).Return(nil)

		err := tool.StoreMetrics(ctx, metricsCounter)
		require.NoError(t, err)
	})

	t.Run("counter with error repo method", func(t *testing.T) {
		tool, repoMock := metricsTool(t)
		ctx := context.Background()

		var delta entity.Counter = 5
		metricsCounter := entity.Metrics{ID: "id", MType: Counter, Delta: &delta}

		var old entity.Counter = 0
		repoMock.On("GetMetrics", ctx, metricsCounter.ID).Return(entity.Metrics{Delta: &old}, repo.ErrNotFound)
		repoMock.On("StoreMetrics", ctx, metricsCounter).Return(errors.New("some method"))

		err := tool.StoreMetrics(ctx, metricsCounter)
		require.Error(t, err)
	})

	t.Run("gauge with error of StoreAll method", func(t *testing.T) {
		tool, repoMock := metricsTool(t)
		tool.syncWriteFile = true
		ctx := context.Background()
		metricsGauge := entity.Metrics{ID: "id", MType: Gauge}
		repoMock.On("StoreMetrics", ctx, metricsGauge).Return(nil)
		repoMock.On("StoreAll", mock.Anything).Return(errors.New("some error"))
		err := tool.StoreMetrics(ctx, metricsGauge)
		require.Error(t, err)
	})
}
