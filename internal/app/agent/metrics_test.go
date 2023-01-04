package agent

import (
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
		collector *storage
	}
	tests := []struct {
		name   string
		fields fields
		want   entity.Counter
	}{
		{
			name: "simple UpdateMetrics test #1",
			fields: fields{
				50,
				&sync.Mutex{},
				&storage{},
			},
			want: 51,
		},
		{
			name: "simple UpdateMetrics test #2",
			fields: fields{
				500,
				&sync.Mutex{},
				&storage{},
			},
			want: 501,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				PollCount: tt.fields.PollCount,
				mu:        tt.fields.Mutex,
				storage:   tt.fields.collector,
			}
			m.UpdateMetrics()
			require.Equal(t, m.PollCount, tt.want)
		})
	}
}

func TestMetrics_collectMetrics(t *testing.T) {
	type fields struct {
		PollCount   entity.Counter
		RandomValue entity.Gauge
		Mutex       *sync.Mutex
		storage     *storage
	}
	type args struct {
		memStats *runtime.MemStats
	}
	type wants struct {
		memStats *runtime.MemStats
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   wants
	}{
		{
			name: "simple collectMetrics test #1",
			fields: fields{
				10,
				10,
				&sync.Mutex{},
				&storage{},
			},
			args: args{
				&runtime.MemStats{Alloc: 1000, Frees: 2000},
			},
			want: wants{
				&runtime.MemStats{Alloc: 1000, Frees: 2000},
			},
		},
		{
			name: "simple collectMetrics test #2",
			fields: fields{
				100,
				100,
				&sync.Mutex{},
				&storage{},
			},
			args: args{
				&runtime.MemStats{TotalAlloc: 2000, HeapAlloc: 4000},
			},
			want: wants{
				&runtime.MemStats{TotalAlloc: 2000, HeapAlloc: 4000},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Metrics{
				PollCount:   tt.fields.PollCount,
				RandomValue: tt.fields.RandomValue,
				mu:          tt.fields.Mutex,
				storage:     tt.fields.storage,
			}
			m.collectMetrics(tt.args.memStats)
			require.Equal(t, m.Alloc, entity.Gauge(tt.want.memStats.Alloc))
			require.Equal(t, m.TotalAlloc, entity.Gauge(tt.want.memStats.TotalAlloc))
			require.Equal(t, m.HeapAlloc, entity.Gauge(tt.want.memStats.HeapAlloc))
			require.Equal(t, m.Frees, entity.Gauge(tt.want.memStats.Frees))
		})
	}
}
