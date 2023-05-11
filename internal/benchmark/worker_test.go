package benchmark

import (
	"github.com/vladislaoramos/alemetric/internal/app/agent"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/internal/usecase"
	"reflect"
	"strings"
	"testing"
)

var metricsNames = []string{
	"Alloc",
	"BuckHashSys",
	"Frees",
	"GCCPUFraction",
	"GCSys",
	"HeapAlloc",
	"HeapIdle",
	"HeapInuse",
	"HeapObjects",
	"HeapReleased",
	"HeapSys",
	"LastGC",
	"Lookups",
	"MCacheInuse",
	"MCacheSys",
	"MSpanInuse",
	"MSpanSys",
	"Mallocs",
	"NextGC",
	"NumForcedGC",
	"NumGC",
	"OtherSys",
	"PauseTotalNs",
	"StackInuse",
	"StackSys",
	"Sys",
	"TotalAlloc",
	"RandomValue",
	"PollCount",
	"TotalMemory",
	"FreeMemory",
	"CPUutilization1",
}

func BenchmarkCollectMetricsStruct(b *testing.B) {
	for i := 0; i < b.N; i++ {
		metrics := agent.NewMetrics()
		tasks := make([]entity.Metrics, 0, len(metricsNames))

		for _, name := range metricsNames {
			field := reflect.Indirect(reflect.ValueOf(metrics)).FieldByName(name)
			if !field.IsValid() {
				continue
			}

			fieldType := strings.ToLower(field.Type().Name())

			var (
				valCounter *entity.Counter
				valGauge   *entity.Gauge
			)

			switch fieldType {
			case usecase.Counter:
				val := entity.Counter(field.Int())
				valCounter = &val
			case usecase.Gauge:
				val := entity.Gauge(field.Float())
				valGauge = &val
			default:
				continue
			}

			task := entity.Metrics{
				ID:    name,
				MType: fieldType,
				Delta: valCounter,
				Value: valGauge,
			}

			tasks = append(tasks, task)
		}
	}
}

func initMetricsStorage() map[string]entity.Metrics {
	metricsStorage := make(map[string]entity.Metrics)

	var (
		gaugeStub   *entity.Gauge
		counterStub *entity.Counter
	)

	slice := []entity.Metrics{
		{
			ID:    "Alloc",
			MType: "BuckHashSys",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "Frees",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "GCCPUFraction",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "GCSys",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "HeapAlloc",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "HeapIdle",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "HeapInuse",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "HeapObjects",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "HeapReleased",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "HeapSys",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "LastGC",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "Lookups",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "MCacheInuse",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "MCacheSys",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "MSpanInuse",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "MSpanSys",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "Mallocs",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "NextGC",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "NumForcedGC",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "NumGC",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "OtherSys",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "PauseTotalNs",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "StackInuse",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "StackSys",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "Sys",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "TotalAlloc",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "TotalMemory",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "FreeMemory",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "CPUutilization1",
			MType: "",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
		{
			ID:    "PollCount",
			MType: "counter",
			Delta: counterStub,
			Value: nil,
			Hash:  "",
		},
		{
			ID:    "RandomValue",
			MType: "gauge",
			Delta: nil,
			Value: gaugeStub,
			Hash:  "",
		},
	}

	for _, item := range slice {
		metricsStorage[item.ID] = item
	}

	return metricsStorage
}

func BenchmarkCollectMetricsMap(b *testing.B) {
	for i := 0; i < b.N; i++ {
		metricsStorage := initMetricsStorage()

		tasks := make([]entity.Metrics, 0, len(metricsNames))

		for _, metrics := range metricsStorage {
			if metrics.MType != usecase.Counter && metrics.MType != usecase.Gauge {
				continue
			}

			if metrics.Delta != nil && metrics.Value != nil {
				continue
			}

			tasks = append(tasks, metrics)
		}
	}
}
