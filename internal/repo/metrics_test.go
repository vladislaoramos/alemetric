package repo

import (
	"github.com/stretchr/testify/require"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"testing"
)

func TestMetricRepo_StoreGauge(t *testing.T) {
	type args struct {
		name  string
		value entity.Gauge
	}

	tests := []struct {
		name   string
		fields map[string]interface{}
		args   args
		want   entity.Gauge
	}{
		{
			name:   "simple test with success",
			fields: map[string]interface{}{},
			args:   args{"TotalAlloc", 100.500},
			want:   entity.Gauge(100.500),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricsRepo{
				storage: tt.fields,
			}

			err := r.StoreGaugeMetrics(tt.args.name, tt.args.value)
			require.NoError(t, err)

			got, ok := r.storage[tt.args.name]
			require.True(t, ok)
			require.Equal(t, got, tt.want)
		})
	}
}

func TestMetricRepo_StoreCounter(t *testing.T) {
	type args struct {
		name  string
		value entity.Counter
	}

	tests := []struct {
		name   string
		fields map[string]interface{}
		args   args
		want   entity.Counter
	}{
		{
			name:   "simple add test with success #1",
			fields: map[string]interface{}{"Total": entity.Counter(100)},
			args:   args{"Total", 1},
			want:   entity.Counter(101),
		},
		{
			name:   "simple add test success #2",
			fields: map[string]interface{}{},
			args:   args{"Total", 1},
			want:   entity.Counter(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &MetricsRepo{
				storage: tt.fields,
			}

			err := r.StoreCounterMetrics(tt.args.name, tt.args.value)
			require.NoError(t, err)

			got, ok := r.storage[tt.args.name]
			require.True(t, ok)
			require.Equal(t, got, tt.want)
		})
	}
}
