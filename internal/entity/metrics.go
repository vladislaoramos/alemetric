// Package entity contains Metrics entity its methods.
// Server responses such metrics to requests of Agent.
package entity

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	grpcTool "github.com/vladislaoramos/alemetric/proto"
	"log"
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

// Metrics stores data of a metrics.
// It contains such attributes: ID, MType, Delta, Value, Hash.
// A correct metrics must have either Delta or Value.
type Metrics struct {
	ID    string   `json:"id" db:"name"`    // metrics name
	MType string   `json:"type" db:"mtype"` // metrics type: either gauge or counter
	Delta *Counter `json:"delta,omitempty"` // metrics value if the type is counter
	Value *Gauge   `json:"value,omitempty"` // metrics value if the type is  gauge
	Hash  string   `json:"hash,omitempty"`  // a hash function value
}

// ParseGaugeMetrics parses Metrics with Gauge type.
func ParseGaugeMetrics(value string) (Gauge, error) {
	s, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return 0, fmt.Errorf("ParseGaugeMetrics - strconv.ParseFloat cannot parse: %w", err)
	}

	return Gauge(s), nil
}

// ParseCounterMetrics parses Metrics with Counter type.
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

// SignData signs an object of metrics.
func (m *Metrics) SignData(app, key string) {
	if key != "" {
		m.Hash = m.hash(key)
		log.Printf("%s: for metric %s made hash %s via key %s\n", app, m.ID, m.Hash, key)
		return
	}
	// log.Printf("%s: signing key is not defined\n", app)
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

// CheckDataSign checks if a metrics is signed.
func (m *Metrics) CheckDataSign(key string) bool {
	return m.Hash == m.hash(key)
}

func (m *Metrics) AsProto() *grpcTool.Metrics {
	var res grpcTool.Metrics

	res.Id = m.ID
	res.Type = m.TypeAsProto()
	res.Hash = m.Hash

	if m.Value != nil {
		res.Payload = &grpcTool.Metrics_Value{Value: float64(*m.Value)}
	} else if m.Delta != nil {
		res.Payload = &grpcTool.Metrics_Delta{Delta: int64(*m.Delta)}
	}

	return &res
}

func (m *Metrics) TypeAsProto() grpcTool.MetricsType {
	switch m.MType {
	case gauge:
		return grpcTool.MetricsType_GAUGE
	case counter:
		return grpcTool.MetricsType_COUNTER
	default:
		return grpcTool.MetricsType_UNKNOWN
	}
}

func FromProto(m *grpcTool.Metrics) Metrics {
	metrics := Metrics{
		ID:    m.GetId(),
		MType: m.GetType().String(),
		Hash:  m.GetHash(),
	}

	if delta := m.GetDelta(); delta != 0 {
		metrics.Delta = (*Counter)(&delta)
	}

	if value := m.GetValue(); value != 0 {
		metrics.Value = (*Gauge)(&value)
	}

	return metrics
}
