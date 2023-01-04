package repo

import (
	"fmt"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"sync"
)

type MetricsRepo struct {
	storage map[string]interface{}
	mu      *sync.Mutex
}

func NewMetricsRepo() *MetricsRepo {
	return &MetricsRepo{
		mu:      &sync.Mutex{},
		storage: make(map[string]interface{}),
	}
}

func (mr *MetricsRepo) StoreGaugeMetrics(name string, value entity.Gauge) error {
	mr.mu.Lock()
	mr.storage[name] = value
	mr.mu.Unlock()
	return nil
}

func (mr *MetricsRepo) StoreCounterMetrics(name string, value entity.Counter) error {
	//mr.mu.Lock()
	if prev, ok := mr.storage[name]; !ok {
		mr.storage[name] = value
	} else {
		mr.storage[name] = value + prev.(entity.Counter)
	}
	//mr.mu.Unlock()

	return nil
}

func (mr *MetricsRepo) GetMetricsNames() []string {
	var ans []string
	for name := range mr.storage {
		ans = append(ans, name)
	}
	return ans
}

func (mr *MetricsRepo) GetMetrics(name string) (interface{}, error) {
	val, ok := mr.storage[name]

	if !ok {
		return nil, fmt.Errorf("%s metrics is not found", name)
	}

	return val, nil
}
