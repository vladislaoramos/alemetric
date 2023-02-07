package repo

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"os"
	"sync"
)

type MetricsRepo struct {
	storage       map[string]entity.Metrics
	Mu            *sync.Mutex
	StoreFilePath string
	Restore       bool
}

func NewMetricsRepo(options ...OptionFunc) *MetricsRepo {
	metricsRepo := &MetricsRepo{
		Mu:      &sync.Mutex{},
		storage: make(map[string]entity.Metrics),
	}

	for _, o := range options {
		o(metricsRepo)
	}

	if metricsRepo.Restore {
		metricsRepo.UploadFromFile()
	}

	return metricsRepo
}

func (r *MetricsRepo) GetMetricsNames() []string {
	var list []string
	for name := range r.storage {
		list = append(list, name)
	}
	return list
}

func (r *MetricsRepo) StoreMetrics(metrics entity.Metrics) error {
	r.Mu.Lock()
	r.storage[metrics.ID] = metrics
	r.Mu.Unlock()
	return nil
}

func (r *MetricsRepo) GetMetrics(name string) (entity.Metrics, error) {
	value, ok := r.storage[name]
	if !ok {
		return entity.Metrics{}, ErrNotFound
	}
	return value, nil
}

func (r *MetricsRepo) StoreToFile() error {
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

func (r *MetricsRepo) UploadFromFile() error {
	file, err := os.OpenFile(r.StoreFilePath, os.O_RDONLY, 0777)
	if err != nil {
		return fmt.Errorf("error opening file with metrics: %w", err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	data, err := reader.ReadBytes('\n')
	if err != nil {
		return fmt.Errorf("error reading from file with metrics: %w", err)
	}

	err = json.Unmarshal(data, &r.storage)
	if err != nil {
		return fmt.Errorf("error unmarshalling file with metrics: %w", err)
	}

	return nil
}
