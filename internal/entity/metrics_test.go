package entity

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseGauge(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Gauge
		wantErr bool
	}{
		{
			name:    "simple test with success",
			args:    "100.500",
			want:    Gauge(100.500),
			wantErr: false,
		},
		{
			name:    "simple test with error",
			args:    "zero",
			want:    Gauge(0),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGaugeMetrics(tt.args)
			if !tt.wantErr {
				require.Equal(t, got, tt.want)
				return
			}
			require.Error(t, err)
		})
	}
}

func TestParseCounter(t *testing.T) {
	tests := []struct {
		name    string
		args    string
		want    Counter
		wantErr bool
	}{
		{
			name:    "simple test with success",
			args:    "1",
			want:    Counter(1),
			wantErr: false,
		},
		{
			name:    "simple test with error",
			args:    "zero",
			want:    Counter(0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCounterMetrics(tt.args)
			if !tt.wantErr {
				require.Equal(t, got, tt.want)
				return
			}
			require.Error(t, err)
		})
	}
}
