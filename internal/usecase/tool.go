package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/internal/repo"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

type ToolUseCase struct {
	repo   MetricsRepo
	logger logger.LogInterface

	writeFileDuration       time.Duration
	writeToFileWithDuration bool
	syncWriteFile           bool
	asyncWriteFile          bool
	C                       chan struct{}

	checkDataSign bool
	encryptionKey string
}

func NewMetricsTool(repo MetricsRepo, l logger.LogInterface, options ...OptionFunc) *ToolUseCase {
	useCase := &ToolUseCase{repo: repo, logger: l}

	for _, o := range options {
		o(useCase)
	}

	if useCase.writeToFileWithDuration {
		go func() {
			ticker := time.NewTicker(useCase.writeFileDuration)
			for {
				<-ticker.C
				useCase.C <- struct{}{}
			}
		}()
	}

	if useCase.writeToFileWithDuration || useCase.asyncWriteFile {
		useCase.C = make(chan struct{}, 1)
		go useCase.saveStorage()
	}

	return useCase
}

func (mt *ToolUseCase) saveStorage() {
	for {
		<-mt.C
		err := mt.repo.StoreAll()
		if err != nil {
			mt.logger.Error(fmt.Sprintf("error while writing to file: %s", err))
		} else {
			mt.logger.Info("store metric success")
		}
	}
}

func (mt *ToolUseCase) GetMetricsNames(ctx context.Context) ([]string, error) {
	names := mt.repo.GetMetricsNames(ctx)
	return names, nil
}

func (mt *ToolUseCase) StoreMetrics(ctx context.Context, metrics entity.Metrics) error {
	if mt.checkDataSign && !metrics.CheckDataSign(mt.encryptionKey) {
		return ErrDataSignNotEqual
	}

	switch metrics.MType {
	case Gauge:
		if err := mt.repo.StoreMetrics(ctx, metrics); err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				return ErrNotFound
			}
			return fmt.Errorf("error store metrics: %w", err)
		}
	case Counter:
		oldMetric, err := mt.repo.GetMetrics(ctx, metrics.ID)
		if err != nil && !errors.Is(err, repo.ErrNotFound) {
			return fmt.Errorf("error getting metrics: %w", err)
		} else if errors.Is(err, repo.ErrNotFound) {
			var oldDelta entity.Counter
			newDelta := oldDelta + *metrics.Delta
			metrics.Delta = &newDelta
		} else {
			delta := *oldMetric.Delta + *metrics.Delta
			metrics.Delta = &delta
		}

		metrics.SignData("server", mt.encryptionKey)

		if err := mt.repo.StoreMetrics(ctx, metrics); err != nil {
			return fmt.Errorf("error storing metrics: %w", err)
		}

	default:
		return ErrNotImplemented
	}
	if mt.asyncWriteFile {
		mt.C <- struct{}{}
	}
	if mt.syncWriteFile {
		err := mt.repo.StoreAll()
		if err != nil {
			return fmt.Errorf("error storing all metrics: %w", err)
		}
	}
	return nil
}

func (mt *ToolUseCase) GetMetrics(ctx context.Context, metrics entity.Metrics) (entity.Metrics, error) {
	res, err := mt.repo.GetMetrics(ctx, metrics.ID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return res, ErrNotFound
		}
		return res, fmt.Errorf("error getting metrics: %w", err)
	}

	return res, nil
}

func (mt *ToolUseCase) PingRepo(ctx context.Context) error {
	return mt.repo.Ping(ctx)
}
