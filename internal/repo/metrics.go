package repo

import (
	"fmt"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"sync"
)

type MetricsRepo struct {
	storage map[string]interface{}
	Mu      *sync.Mutex
}

func NewMetricsRepo() *MetricsRepo {
	return &MetricsRepo{
		Mu:      &sync.Mutex{},
		storage: make(map[string]interface{}),
	}
}

func (r *MetricsRepo) GetMetricsNames() []string {
	var list []string
	for name := range r.storage {
		list = append(list, name)
	}
	return list
}

func (r *MetricsRepo) StoreGaugeMetrics(name string, value entity.Gauge) error {
	r.Mu.Lock()
	r.storage[name] = value
	r.Mu.Unlock()
	return nil
}

func (r *MetricsRepo) StoreCounterMetrics(name string, value entity.Counter) error {
	r.Mu.Lock()
	oldValue, ok := r.storage[name]
	r.Mu.Unlock()
	if ok {
		r.storage[name] = value + oldValue.(entity.Counter)
	} else {
		r.storage[name] = value
	}
	return nil
}

func (r *MetricsRepo) GetMetrics(name string) (interface{}, error) {
	value, ok := r.storage[name]
	if !ok {
		return nil, fmt.Errorf("%s metrics is not found", name)
	}
	return value, nil
}
