package entity

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseGaugeMetrics(t *testing.T) {
	type args struct {
		value string
	}

	tests := []struct {
		name    string
		args    args
		want    Gauge
		wantErr bool
	}{
		{
			name:    "simple ParseGauge with success",
			args:    args{"1.1"},
			want:    Gauge(1.1),
			wantErr: false,
		},
		{
			name:    "simple ParseGauge with error",
			args:    args{"zero"},
			want:    Gauge(0),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseGaugeMetrics(tt.args.value)
			if !tt.wantErr {
				require.Equal(t, got, tt.want)
				return
			}
			require.Error(t, err)
		})
	}
}

func TestParseCounterMetrics(t *testing.T) {
	type args struct {
		value string
	}
	tests := []struct {
		name    string
		args    args
		want    Counter
		wantErr bool
	}{
		{
			name:    "simple ParseCounter with success",
			args:    args{"1"},
			want:    Counter(1),
			wantErr: false,
		},
		{
			name:    "simple ParseGauge with error",
			args:    args{"zero"},
			want:    Counter(0),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseCounterMetrics(tt.args.value)
			if !tt.wantErr {
				require.Equal(t, got, tt.want)
				return
			}
			require.Error(t, err)
		})
	}
}

func TestCheckDataSign(t *testing.T) {
	key := "secretKey"
	var value Gauge = 1.23
	metrics := &Metrics{
		ID:    "metric1",
		MType: "gauge",
		Value: &value,
		Hash:  "",
	}

	expectedHash := metrics.hash(key)
	metrics.Hash = expectedHash
	result := metrics.CheckDataSign(key)

	require.True(t, result)
}

func TestSignData(t *testing.T) {
	key := "secretKey"
	var value Gauge = 1.23
	metrics := &Metrics{
		ID:    "metric1",
		MType: "gauge",
		Value: &value,
		Hash:  "",
	}

	metrics.SignData("TestApp", key)
	require.NotEqual(t, "", metrics.Hash)
	metrics.SignData("TestApp", "")
	require.NotEqual(t, "", metrics.Hash)
}
