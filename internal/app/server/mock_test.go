package server

import "github.com/vladislaoramos/alemetric/internal/entity"

type MockMetricsRepo struct {
	MockMetrics interface{}
	MockErr     error
}

func (m *MockMetricsRepo) StoreGaugeMetrics(string, entity.Gauge) error     { return m.MockErr }
func (m *MockMetricsRepo) StoreCounterMetrics(string, entity.Counter) error { return m.MockErr }
func (m *MockMetricsRepo) GetMetrics(string) (interface{}, error)           { return m.MockMetrics, m.MockErr }
func (m *MockMetricsRepo) GetMetricsNames() []string                        { return nil }
