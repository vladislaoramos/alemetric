package entity

import (
	"fmt"
	"strconv"
)

type (
	Gauge   float64
	Counter int64
)

type Metrics struct {
	ID    string   `json:"id"`              // имя метрики
	MType string   `json:"type"`            // параметр, принимающий значение gauge или counter
	Delta *Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
}

func ParseGaugeMetrics(value string) (Gauge, error) {
	s, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("error with ParseGaugeMetrics: strconv.ParseFloat cannot parse: %w", err)
	}

	return Gauge(s), nil
}

func ParseCounterMetrics(value string) (Counter, error) {
	s, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("error with ParseCounterMetrics strconv.Atoi cannot parse: %w", err)
	}
	return Counter(s), nil
}

func (g Gauge) Type() string {
	return "gauge"
}

func (c Counter) Type() string {
	return "counter"
}
