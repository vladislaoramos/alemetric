package server

import (
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
