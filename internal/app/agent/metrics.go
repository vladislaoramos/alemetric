package agent

import (
	"github.com/vladislaoramos/alemetric/internal/entity"
	"math/rand"
	"runtime"
	"sync"
)

type Metrics struct {
	PollCount   entity.Counter
	RandomValue entity.Gauge
	Mu          *sync.Mutex
	*storage
}

type storage struct {
	Alloc         entity.Gauge
	BuckHashSys   entity.Gauge
	Frees         entity.Gauge
	GCCPUFraction entity.Gauge
	GCSys         entity.Gauge
	HeapAlloc     entity.Gauge
	HeapIdle      entity.Gauge
	HeapInuse     entity.Gauge
	HeapObjects   entity.Gauge
	HeapReleased  entity.Gauge
	HeapSys       entity.Gauge
	LastGC        entity.Gauge
	Lookups       entity.Gauge
	MCacheInuse   entity.Gauge
	MCacheSys     entity.Gauge
	MSpanInuse    entity.Gauge
	MSpanSys      entity.Gauge
	Mallocs       entity.Gauge
	NextGC        entity.Gauge
	NumForcedGC   entity.Gauge
	NumGC         entity.Gauge
	OtherSys      entity.Gauge
	PauseTotalNs  entity.Gauge
	StackInuse    entity.Gauge
	StackSys      entity.Gauge
	Sys           entity.Gauge
	TotalAlloc    entity.Gauge
}

func NewMetrics() *Metrics {
	return &Metrics{
		Mu:      &sync.Mutex{},
		storage: &storage{},
	}
}

func (m *Metrics) CollectMetrics() {
	memStats := &runtime.MemStats{}
	runtime.ReadMemStats(memStats)

	m.Mu.Lock()
	defer m.Mu.Unlock()

	m.updateMetrics(memStats)
	m.PollCount += 1
	m.RandomValue = entity.Gauge(rand.Float64())
	// m.Mu.Unlock()
}

func (m *Metrics) updateMetrics(memStats *runtime.MemStats) {
	m.storage.Alloc = entity.Gauge(memStats.Alloc)
	m.storage.BuckHashSys = entity.Gauge(memStats.BuckHashSys)
	m.storage.Frees = entity.Gauge(memStats.Frees)
	m.storage.GCCPUFraction = entity.Gauge(memStats.GCCPUFraction)
	m.storage.GCSys = entity.Gauge(memStats.GCSys)
	m.storage.HeapAlloc = entity.Gauge(memStats.HeapAlloc)
	m.storage.HeapIdle = entity.Gauge(memStats.HeapIdle)
	m.storage.HeapInuse = entity.Gauge(memStats.HeapInuse)
	m.storage.HeapObjects = entity.Gauge(memStats.HeapObjects)
	m.storage.HeapReleased = entity.Gauge(memStats.HeapReleased)
	m.storage.HeapSys = entity.Gauge(memStats.HeapSys)
	m.storage.LastGC = entity.Gauge(memStats.LastGC)
	m.storage.Lookups = entity.Gauge(memStats.Lookups)
	m.storage.MCacheInuse = entity.Gauge(memStats.MCacheInuse)
	m.storage.MCacheSys = entity.Gauge(memStats.MCacheSys)
	m.storage.MSpanInuse = entity.Gauge(memStats.MSpanInuse)
	m.storage.MSpanSys = entity.Gauge(memStats.MSpanSys)
	m.storage.Mallocs = entity.Gauge(memStats.Mallocs)
	m.storage.NextGC = entity.Gauge(memStats.NextGC)
	m.storage.NumForcedGC = entity.Gauge(memStats.NumForcedGC)
	m.storage.NumGC = entity.Gauge(memStats.NumGC)
	m.storage.OtherSys = entity.Gauge(memStats.OtherSys)
	m.storage.PauseTotalNs = entity.Gauge(memStats.PauseTotalNs)
	m.storage.StackInuse = entity.Gauge(memStats.StackInuse)
	m.storage.StackSys = entity.Gauge(memStats.StackSys)
	m.storage.Sys = entity.Gauge(memStats.Sys)
	m.storage.TotalAlloc = entity.Gauge(memStats.TotalAlloc)
}
