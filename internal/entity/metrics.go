package entity

import (
	"strconv"
)

type (
	Gauge   float64
	Counter int64
)

type Metrics struct {
	Name  string
	Value interface{}
}

func ParseGaugeMetrics(value string) (Gauge, error) {
	s, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, err
	}
	return Gauge(s), nil
}

func ParseCounterMetrics(value string) (Counter, error) {
	s, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return Counter(s), nil
}
