package repo

import (
	"github.com/vladislaoramos/alemetric/internal/entity"
)

type MetricsRepo struct {
	storage map[string]interface{}
}

func NewMetricsRepo() *MetricsRepo {
	mr := new(MetricsRepo)
	mr.storage = make(map[string]interface{})
	return mr
}

func (mr *MetricsRepo) StoreGaugeMetrics(name string, value entity.Gauge) error {
	mr.storage[name] = value
	return nil
}

func (mr *MetricsRepo) StoreCounterMetrics(name string, value entity.Counter) error {
	if prev, ok := mr.storage[name]; !ok {
		mr.storage[name] = value
	} else {
		mr.storage[name] = value + prev.(entity.Counter)
	}

	return nil
}
