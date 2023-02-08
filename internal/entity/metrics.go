package entity

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"strconv"
)

type (
	Gauge   float64
	Counter int64
)

const (
	counter = "counter"
	gauge   = "gauge"
)

type Metrics struct {
	ID    string   `json:"id" db:"name"`    // имя метрики
	MType string   `json:"type" db:"mtype"` // параметр, принимающий значение gauge или counter
	Delta *Counter `json:"delta,omitempty"` // значение метрики в случае передачи counter
	Value *Gauge   `json:"value,omitempty"` // значение метрики в случае передачи gauge
	Hash  string   `json:"hash,omitempty"`  // значение хеш-функции
}

func ParseGaugeMetrics(value string) (Gauge, error) {
	s, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("ParseGaugeMetrics - strconv.ParseFloat cannot parse: %w", err)
	}

	return Gauge(s), nil
}

func ParseCounterMetrics(value string) (Counter, error) {
	s, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("ParseCounterMetrics - strconv.Atoi cannot parse: %w", err)
	}
	return Counter(s), nil
}

func (g Gauge) Type() string {
	return gauge
}

func (c Counter) Type() string {
	return counter
}

func (m *Metrics) SignData(key string) {
	if key != "" {
		m.Hash = m.hash(key)
	}
}

func (m *Metrics) hash(key string) string {
	var res string

	switch m.MType {
	case gauge:
		res = fmt.Sprintf("%s:%s:%f", m.ID, m.MType, *m.Value)
	case counter:
		res = fmt.Sprintf("%s:%s:%d", m.ID, m.MType, *m.Delta)
	}

	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(res))

	return fmt.Sprintf("%x", h.Sum(nil))
}

func (m *Metrics) CheckDataSign(key string) bool {
	return m.Hash == m.hash(key)
}
