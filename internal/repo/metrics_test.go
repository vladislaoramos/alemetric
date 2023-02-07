package repo

//import (
//	"github.com/stretchr/testify/require"
//	"github.com/vladislaoramos/alemetric/internal/entity"
//	"sync"
//	"testing"
//)
//
//func TestMetricRepo_StoreGauge(t *testing.T) {
//	type args struct {
//		name  string
//		value entity.Gauge
//	}
//
//	tests := []struct {
//		name   string
//		fields map[string]interface{}
//		args   args
//		want   entity.Gauge
//	}{
//		{
//			name:   "simple test with TotalAlloc",
//			fields: map[string]interface{}{},
//			args:   args{"TotalAlloc", 100.500},
//			want:   entity.Gauge(100.500),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &MetricsRepo{
//				storage: tt.fields,
//				Mu:      &sync.Mutex{},
//			}
//			err := r.StoreGaugeMetrics(tt.args.name, tt.args.value)
//			require.NoError(t, err)
//
//			got, ok := r.storage[tt.args.name]
//			require.True(t, ok)
//			require.Equal(t, got, tt.want)
//		})
//	}
//}
//
//func TestMetricRepo_AddCounter(t *testing.T) {
//	type args struct {
//		name  string
//		value entity.Counter
//	}
//
//	tests := []struct {
//		name   string
//		fields map[string]interface{}
//		args   args
//		want   entity.Counter
//	}{
//		{
//			name:   "simple test with PollCount #1",
//			fields: map[string]interface{}{},
//			args:   args{"PollCount", 1},
//			want:   entity.Counter(1),
//		},
//		{
//			name:   "simple test with PollCount #1",
//			fields: map[string]interface{}{"PollCount": entity.Counter(1)},
//			args:   args{"PollCount", 1},
//			want:   entity.Counter(2),
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			r := &MetricsRepo{
//				storage: tt.fields,
//				Mu:      &sync.Mutex{},
//			}
//
//			err := r.StoreCounterMetrics(tt.args.name, tt.args.value)
//			require.NoError(t, err)
//
//			got, ok := r.storage[tt.args.name]
//			require.True(t, ok)
//			require.Equal(t, got, tt.want)
//		})
//	}
//}
