package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func TestRouter(t *testing.T) {
	type want struct {
		code int
		body string
	}

	tests := []struct {
		name    string
		args    MetricsRepo
		request string
		method  string
		want    want
	}{
		{
			name:    "simple test of main request with success",
			args:    &MockMetricsRepo{},
			request: "/",
			method:  http.MethodGet,
			want: want{
				code: 200,
				body: "",
			},
		},
		{
			name:    "simple test of updating gauge value with success",
			args:    &MockMetricsRepo{},
			request: "/update/gauge/TotalAlloc/1.1",
			method:  http.MethodPost,
			want: want{
				code: 200,
				body: "",
			},
		},
		{
			name:    "simple test of updating counter value with success",
			args:    &MockMetricsRepo{},
			request: "/update/counter/Frees/1",
			method:  http.MethodPost,
			want: want{
				code: 200,
				body: "",
			},
		},
		{
			name:    "simple test of getting gauge value with success",
			args:    &MockMetricsRepo{MockMetrics: entity.Gauge(777)},
			request: "/value/gauge/NumGC",
			method:  http.MethodGet,
			want: want{
				code: 200,
				body: "777",
			},
		},
		{
			name:    "simple test of getting counter value with success",
			args:    &MockMetricsRepo{MockMetrics: entity.Counter(777)},
			request: "/value/counter/NextGC",
			method:  http.MethodGet,
			want: want{
				code: 200,
				body: "777",
			},
		},
		{
			name:    "simple test of not implemented method",
			args:    &MockMetricsRepo{},
			request: "/update/range/NextGC/1.1",
			method:  http.MethodPost,
			want: want{
				code: 501,
				body: "metrics type is not found\n",
			},
		},
		{
			name:    "simple test of updating gauge value with error",
			args:    &MockMetricsRepo{MockErr: errors.New("error")},
			request: "/update/gauge/NextGC/1.1",
			method:  http.MethodPost,
			want: want{
				code: 500,
				body: "some problem with storage\n",
			},
		},
		{
			name:    "simple test of getting gauge value with error",
			args:    &MockMetricsRepo{MockErr: errors.New("error")},
			request: "/value/gauge/Frees",
			method:  http.MethodGet,
			want: want{
				code: 404,
				body: "metrics is not found\n",
			},
		},
	}
	for _, tt := range tests {
		r := chi.NewRouter()
		NewRouter(r, tt.args)
		ts := httptest.NewServer(r)

		req, err := http.NewRequest(tt.method, ts.URL+tt.request, nil)
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		respBody, err := io.ReadAll(resp.Body)
		require.NoError(t, err)

		body := string(respBody)

		resp.Body.Close()
		ts.Close()

		require.Equal(t, tt.want.code, resp.StatusCode)
		require.Equal(t, tt.want.body, body)
	}
}

func TestRouterJSON(t *testing.T) {
	type args struct {
		repo    MetricsRepo
		metrics entity.Metrics
	}

	type want struct {
		code  int
		value *entity.Gauge
		delta *entity.Counter
	}

	var (
		gaugeVal   = entity.Gauge(2.4)
		counterVal = entity.Counter(2)
	)

	tests := []struct {
		name    string
		args    args
		request string
		method  string
		want    want
	}{
		{
			name: "simple test of updating gauge value with success",
			args: args{
				&MockMetricsRepo{},
				entity.Metrics{
					ID:    "TotalAlloc",
					MType: "gauge",
					Delta: nil,
					Value: &gaugeVal,
				},
			},
			method:  http.MethodPost,
			request: "/update/",
			want: want{
				code:  200,
				value: nil,
				delta: nil,
			},
		},
		{
			name: "simple test of updating counter value with success",
			args: args{
				&MockMetricsRepo{},
				entity.Metrics{
					ID:    "NextGC",
					MType: "counter",
					Delta: &counterVal,
					Value: nil},
			},
			method:  http.MethodPost,
			request: "/update/",
			want: want{
				code:  200,
				value: nil,
				delta: nil,
			},
		},
		{
			name: "simple test of getting gauge value with success",
			args: args{
				&MockMetricsRepo{MockMetrics: gaugeVal},
				entity.Metrics{
					ID:    "TotalAlloc",
					MType: "gauge",
					Delta: nil,
					Value: nil,
				},
			},
			method:  http.MethodPost,
			request: "/value/",
			want: want{
				code:  200,
				value: &gaugeVal,
				delta: nil,
			},
		},
		{
			name: "simple test of getting gauge value with error",
			args: args{
				&MockMetricsRepo{MockErr: errors.New("error")},
				entity.Metrics{
					ID:    "TotalAlloc",
					MType: "gauge",
					Delta: nil,
					Value: nil,
				},
			},
			method:  http.MethodPost,
			request: "/value/",
			want: want{
				code:  404,
				value: nil,
				delta: nil,
			},
		},
	}

	for _, tt := range tests {
		r := chi.NewRouter()
		NewRouter(r, tt.args.repo)
		ts := httptest.NewServer(r)

		reqJSON, err := json.Marshal(tt.args.metrics)
		require.NoError(t, err)

		req, err := http.NewRequest(tt.method, ts.URL+tt.request, bytes.NewBuffer(reqJSON))
		require.NoError(t, err)

		resp, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		require.Equal(t, tt.want.code, resp.StatusCode)

		if tt.want.value != nil || tt.want.delta != nil {
			var respJSON entity.Metrics

			err = json.NewDecoder(resp.Body).Decode(&respJSON)
			require.NoError(t, err)

			require.Equal(t, respJSON.Value, tt.want.value)
			require.Equal(t, respJSON.Delta, tt.want.delta)
		}

		ts.Close()
		resp.Body.Close()
	}
}
