package agent

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

func TestNewWebAPI(t *testing.T) {
	tests := []struct {
		name string
		args *resty.Client
		want *WebAPIClient
	}{
		{
			name: "simple test #1",
			args: &resty.Client{},
			want: &WebAPIClient{client: &resty.Client{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewWebAPI(tt.args)
			require.Equal(t, got, tt.want)
		})
	}
}

func TestWebAPI_SendMetric(t *testing.T) {
	type args struct {
		metricName  string
		metricType  string
		metricValue interface{}
	}

	tests := []struct {
		name           string
		args           args
		response       int
		withTestServer bool
		wantErr        bool
	}{
		{
			name:           "simple test with success",
			args:           args{"Frees", "gauge", 100.500},
			response:       http.StatusOK,
			withTestServer: true,
			wantErr:        false,
		},
		{
			name:           "simple test with error",
			args:           args{"Frees", "gauge", 100.500},
			response:       http.StatusBadRequest,
			withTestServer: true,
			wantErr:        true,
		},
		{
			name:           "simple test without server",
			args:           args{"Frees", "gauge", 100.500},
			withTestServer: false,
			wantErr:        true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.response)
			}))
			defer testServer.Close()

			var serverURL string
			if tt.withTestServer {
				serverURL = testServer.URL
			}

			webAPI := &WebAPIClient{resty.New().SetBaseURL(serverURL)}

			err := webAPI.SendMetrics(tt.args.metricName, tt.args.metricType, tt.args.metricValue)
			if !tt.wantErr {
				return
			}
			require.Error(t, err)
		})
	}
}
