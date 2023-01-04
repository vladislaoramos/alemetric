package agent

import (
	"math/rand"
	"runtime"
	"sync"
)

type (
	gauge   float64
	counter int64
)

type Metrics struct {
	PollCount   counter
	RandomValue gauge
	mu          *sync.Mutex
	*storage
}

type storage struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge
}

func NewMetrics() *Metrics {
	return &Metrics{
		mu:      &sync.Mutex{},
		storage: &storage{},
	}
}

func (m *Metrics) UpdateMetrics() {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	m.mu.Lock()
	// defer m.mu.Unlock()
	m.collectMetrics(memStats)
	m.mu.Unlock()

	m.RandomValue = gauge(rand.NormFloat64())
	m.PollCount += 1
}

func (m *Metrics) collectMetrics(memStats *runtime.MemStats) {
	m.storage.Alloc = gauge(memStats.Alloc)
	m.storage.BuckHashSys = gauge(memStats.BuckHashSys)
	m.storage.Frees = gauge(memStats.Frees)
	m.storage.GCCPUFraction = gauge(memStats.GCCPUFraction)
	m.storage.GCSys = gauge(memStats.GCSys)
	m.storage.HeapAlloc = gauge(memStats.HeapAlloc)
	m.storage.HeapIdle = gauge(memStats.HeapIdle)
	m.storage.HeapInuse = gauge(memStats.HeapInuse)
	m.storage.HeapObjects = gauge(memStats.HeapObjects)
	m.storage.HeapReleased = gauge(memStats.HeapReleased)
	m.storage.HeapSys = gauge(memStats.HeapSys)
	m.storage.LastGC = gauge(memStats.LastGC)
	m.storage.Lookups = gauge(memStats.Lookups)
	m.storage.MCacheInuse = gauge(memStats.MCacheInuse)
	m.storage.MCacheSys = gauge(memStats.MCacheSys)
	m.storage.MSpanInuse = gauge(memStats.MSpanInuse)
	m.storage.MSpanSys = gauge(memStats.MSpanSys)
	m.storage.Mallocs = gauge(memStats.Mallocs)
	m.storage.NextGC = gauge(memStats.NextGC)
	m.storage.NumForcedGC = gauge(memStats.NumForcedGC)
	m.storage.NumGC = gauge(memStats.NumGC)
	m.storage.OtherSys = gauge(memStats.OtherSys)
	m.storage.PauseTotalNs = gauge(memStats.PauseTotalNs)
	m.storage.StackInuse = gauge(memStats.StackInuse)
	m.storage.StackSys = gauge(memStats.StackSys)
	m.storage.Sys = gauge(memStats.Sys)
	m.storage.TotalAlloc = gauge(memStats.TotalAlloc)
}
