package repo

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/vladislaoramos/alemetric/internal/entity"
)

// MetricsRepo stores the object for interaction with the in-memory storage.
type MetricsRepo struct {
	storage       map[string]entity.Metrics
	Mu            *sync.Mutex
	StoreFilePath string
	Restore       bool
}

// NewMetricsRepo creates the in-memory storage object.
func NewMetricsRepo(options ...OptionFunc) (*MetricsRepo, error) {
	metricsRepo := &MetricsRepo{
		Mu:      &sync.Mutex{},
		storage: make(map[string]entity.Metrics),
	}

	for _, o := range options {
		o(metricsRepo)
	}

	if metricsRepo.Restore {
		err := metricsRepo.Upload(context.TODO())
		if err != nil {
			return nil, err
		}
	}

	return metricsRepo, nil
}

// GetMetricsNames gets all metrics names from the in-memory storage.
func (r *MetricsRepo) GetMetricsNames(_ context.Context) []string {
	var list []string
	r.Mu.Lock()
	defer r.Mu.Unlock()
	for name := range r.storage {
		list = append(list, name)
	}
	return list
}

// StoreMetrics stores a metrics into the in-memory storage.
func (r *MetricsRepo) StoreMetrics(_ context.Context, metrics entity.Metrics) error {
	r.Mu.Lock()
	r.storage[metrics.ID] = metrics
	r.Mu.Unlock()
	return nil
}

// GetMetrics gets a metrics from the in-memory storage.
func (r *MetricsRepo) GetMetrics(_ context.Context, name string) (entity.Metrics, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	value, ok := r.storage[name]
	if !ok {
		return entity.Metrics{}, ErrNotFound
	}
	return value, nil
}

// StoreAll stores all metrics from the in-memory storage to the store file.
func (r *MetricsRepo) StoreAll() error {
	file, err := os.OpenFile(r.StoreFilePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("error opening file with metrics: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	r.Mu.Lock()
	data, err := json.Marshal(r.storage)
	r.Mu.Unlock()
	if err != nil {
		return fmt.Errorf("error marshalling file with metrics: %w", err)
	}

	_, err = writer.Write(data)
	if err != nil {
		return fmt.Errorf("error writing to structure with metrics: %w", err)
	}

	if err := writer.WriteByte('\n'); err != nil {
		return fmt.Errorf("error writing to file with metrics: %w", err)
	}

	writer.Flush()

	return nil
}

// Upload uploads all metrics from the store file into the in-memory storage.
func (r *MetricsRepo) Upload(_ context.Context) error {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	file, err := os.OpenFile(r.StoreFilePath, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("error opening file with metrics: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	data, err := reader.ReadBytes('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("error reading from file with metrics: %w", err)
	}
	if errors.Is(err, io.EOF) {
		return nil
	}

	err = json.Unmarshal(data, &r.storage)
	if err != nil {
		return fmt.Errorf("error unmarshalling file with metrics: %w", err)
	}

	return nil
}

func (r *MetricsRepo) Ping(_ context.Context) error {
	return nil
}
