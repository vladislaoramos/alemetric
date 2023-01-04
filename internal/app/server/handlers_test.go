package server

import (
	"github.com/vladislaoramos/alemetric/internal/entity"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type MetricsRepoMock struct{}

func (m *MetricsRepoMock) StoreGaugeMetrics(string, entity.Gauge) error     { return nil }
func (m *MetricsRepoMock) StoreCounterMetrics(string, entity.Counter) error { return nil }

func TestUpdateMetricsHandler(t *testing.T) {
	type args struct {
		repo MetricsRepo
	}

	type test struct {
		name    string
		args    args
		want    int
		request string
		method  string
	}

	tests := []test{
		{
			name: "simple test with success",
			args: args{
				repo: &MetricsRepoMock{},
			},
			want:    http.StatusOK,
			request: "/update/gauge/StackSys/5.8",
			method:  http.MethodPost,
		},
		{
			name: "simple test with method error",
			args: args{
				repo: &MetricsRepoMock{},
			},
			want:    http.StatusMethodNotAllowed,
			request: "/update/gauge/Alloc/1.1",
			method:  http.MethodGet,
		},
		{
			name: "simple test with not found",
			args: args{
				repo: &MetricsRepoMock{},
			},
			want:    http.StatusNotFound,
			request: "/update/gauge/",
			method:  http.MethodPost,
		},
		{
			name: "simple test with bad request",
			args: args{
				repo: &MetricsRepoMock{},
			},
			want:    http.StatusBadRequest,
			request: "/update/counter/testCounter/none",
			method:  http.MethodPost,
		},
		{
			name: "simple test with not implemented",
			args: args{
				repo: &MetricsRepoMock{},
			},
			want:    http.StatusNotImplemented,
			request: "/update/storage/type/some",
			method:  http.MethodPost,
		},
		{
			name: "simple test without metrics value",
			args: args{
				repo: &MetricsRepoMock{},
			},
			want:    http.StatusNotFound,
			request: "/update/gauge/TotalAlloc",
			method:  http.MethodPost,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			request := httptest.NewRequest(tt.method, tt.request, nil)

			response := httptest.NewRecorder()
			h := UpdateMetricsHandler(tt.args.repo)
			h.ServeHTTP(response, request)
			res := response.Result()

			defer res.Body.Close()

			require.Equal(t, tt.want, res.StatusCode)
		})
	}
}
