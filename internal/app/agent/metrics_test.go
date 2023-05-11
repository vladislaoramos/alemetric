package agent

import (
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMetrics_UpdateMetrics(t *testing.T) {
	type fields struct {
		PollCount entity.Counter
		Mutex     *sync.Mutex
		*storage
	}
	tests := []struct {
		name   string
		fields fields
		want   entity.Counter
	}{
		{
			name: "simple test of UpdateMetrics",
			fields: fields{
				100,
				&sync.Mutex{},
				&storage{},
			},
			want: 101,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				PollCount: tt.fields.PollCount,
				Mu:        tt.fields.Mutex,
				storage:   tt.fields.storage,
			}

			m.CollectMetrics()

			require.Equal(t, m.PollCount, tt.want)
		})
	}
}

func TestMetrics_CollectMetrics(t *testing.T) {
	type fields struct {
		PollCount entity.Counter
		Mutex     *sync.Mutex
		collector *storage
	}

	tests := []struct {
		name   string
		fields fields
		args   *runtime.MemStats
		want   entity.Gauge
	}{
		{
			name: "simple test of collectMetrics",
			fields: fields{
				100,
				&sync.Mutex{},
				&storage{},
			},
			args: &runtime.MemStats{Alloc: 1000},
			want: entity.Gauge(1000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				PollCount: tt.fields.PollCount,
				Mu:        tt.fields.Mutex,
				storage:   tt.fields.collector,
			}
			m.updateMetrics(tt.args)
			require.Equal(t, m.Alloc, tt.want)
		})
	}
}

func TestMetrics_CollectAdditionalMetrics(t *testing.T) {
	m := &Metrics{
		PollCount: 100,
		Mu:        &sync.Mutex{},
		storage:   &storage{},
	}

	m.CollectAdditionalMetrics()
	vm, _ := mem.VirtualMemory()

	require.NotNil(t, entity.Gauge(vm.Total), m.TotalMemory)
	require.NotNil(t, entity.Gauge(vm.Free), m.FreeMemory)
}

func TestNewMetrics(t *testing.T) {
	expected := &Metrics{
		Mu:      &sync.Mutex{},
		storage: &storage{},
	}

	actual := NewMetrics()

	require.EqualValues(t, expected, actual)
}
