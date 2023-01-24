package usecase

import (
	"errors"
	"fmt"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/internal/repo"
	"github.com/vladislaoramos/alemetric/pkg/log"
	"time"
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
	C                       chan struct{}
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

	if useCase.writeToFileWithDuration || useCase.syncWriteFile {
		useCase.C = make(chan struct{}, 1)
		go useCase.saveStorage()
	}

	return useCase
}

func (mt *ToolUseCase) saveStorage() {
	for {
		<-mt.C
		err := mt.repo.StoreToFile()
		if err != nil {
			mt.logger.Error(fmt.Sprintf("error while writing to file: %s", err))
		} else {
			mt.logger.Info("store metric success")
		}
	}
}

func (mt *ToolUseCase) GetMetricsNames() ([]string, error) {
	names := mt.repo.GetMetricsNames()
	return names, nil
}

func (mt *ToolUseCase) StoreMetrics(metrics entity.Metrics) error {
	switch metrics.MType {
	case Gauge:
		if err := mt.repo.StoreMetrics(metrics); err != nil {
			if errors.Is(err, repo.ErrNotFound) {
				return ErrNotFound
			}
			return fmt.Errorf("MetricsTool - StoreMetric: %w", err)
		}
	case Counter:
		oldMetric, err := mt.repo.GetMetrics(metrics.ID)
		if err != nil {
			if !errors.Is(err, repo.ErrNotFound) {
				return fmt.Errorf("MetricsTool - GetMetric: %w", err)
			}
		} else {
			delta := *metrics.Delta + *oldMetric.Delta
			metrics.Delta = &delta
		}
		if err := mt.repo.StoreMetrics(metrics); err != nil {
			return fmt.Errorf("MetricsTool - StoreMetric: %w", err)
		}
	default:
		return ErrNotImplemented
	}
	if mt.syncWriteFile {
		mt.C <- struct{}{}
	}
	return nil
}

func (mt *ToolUseCase) GetMetrics(metrics entity.Metrics) (entity.Metrics, error) {
	metric, err := mt.repo.GetMetrics(metrics.ID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			return metric, ErrNotFound
		}
		return metric, fmt.Errorf("MetricsTool - Metric: %w", err)
	}
	return metric, nil
}
