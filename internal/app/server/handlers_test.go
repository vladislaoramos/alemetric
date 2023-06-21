package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/internal/repo"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/vladislaoramos/alemetric/internal/usecase"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

type TestServer struct {
	Server *httptest.Server
}

func NewTestServer(metricsStorage *repo.MetricsRepo, lgr *logger.Logger) TestServer {
	handler := chi.NewRouter()
	mtOptions := make([]usecase.OptionFunc, 0)
	mt := usecase.NewMetricsTool(metricsStorage, lgr, mtOptions...)
	NewRouter(handler, mt, lgr, "", "127.0.0.0/8")
	ts := httptest.NewServer(handler)
	return TestServer{
		Server: ts,
	}
}

func testLogger() *logger.Logger {
	f, err := os.OpenFile("/tmp/test_log_server", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	l := logger.New("debug", f)
	if err != nil {
		l.Fatal("unable to open file for log")
	}

	return l
}

func (s *TestServer) testRequest(
	t *testing.T,
	method, path string,
	body io.Reader) (
	int, []byte,
) {
	req, err := http.NewRequest(method, s.Server.URL+path, body)
	require.NoError(t, err)

	//if body != nil {
	//	req.Header.Add("Content-Type", "application/json")
	//}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Real-IP", "127.0.0.1")

	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)

	respBody, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	defer resp.Body.Close()

	return resp.StatusCode, respBody
}

func TestGetMetricsHandler(t *testing.T) {
	memStorage, err := repo.NewMetricsRepo()
	assert.NoError(t, err)

	tl := testLogger()
	ts := NewTestServer(memStorage, tl)
	for _, request := range []string{
		"/update/gauge/HeapInuse/786432.01",
		"/update/gauge/HeapObjects/613",
		"/update/gauge/NextGC/4194304",
		"/update/gauge/LastGC/0",
	} {
		statusCode, _ := ts.testRequest(t, "POST", request, nil)
		assert.Equal(t, http.StatusOK, statusCode)
	}

	statusCode, _ := ts.testRequest(t, "GET", "/", nil)
	assert.Equal(t, http.StatusOK, statusCode)

	for _, item := range []string{
		"HeapInuse",
		"HeapObjects",
		"NextGC",
		"LastGC",
	} {
		statusCode, _ = ts.testRequest(t, "GET", "/value/gauge/"+item, nil)
		assert.Equal(t, http.StatusOK, statusCode)
	}
}

func TestUpdateSeveralMetricsHandler(t *testing.T) {
	memStorage, err := repo.NewMetricsRepo()
	assert.NoError(t, err)

	tl := testLogger()
	ts := NewTestServer(memStorage, tl)

	var (
		tCounter entity.Counter = 5
		tGauge   entity.Gauge   = 123.01
	)

	type request struct {
		metricsType        string
		metricsName        string
		metricsValue       *entity.Gauge
		metricsDelta       *entity.Counter
		expectedStatusCode int
	}

	requests := []request{
		{
			metricsType:        Counter,
			metricsName:        "PollCount",
			metricsDelta:       &tCounter,
			expectedStatusCode: http.StatusOK,
		},
		{
			metricsType:        Gauge,
			metricsName:        "BuckHashSys",
			metricsValue:       &tGauge,
			expectedStatusCode: http.StatusOK,
		},
	}

	var severalUpdates []entity.Metrics
	for _, req := range requests {
		severalUpdates = append(severalUpdates, entity.Metrics{
			ID:    req.metricsName,
			MType: req.metricsType,
			Delta: req.metricsDelta,
			Value: req.metricsValue,
		})
	}

	b, err := json.Marshal(severalUpdates)
	assert.NoError(t, err)

	reader := strings.NewReader(string(b))

	statusCode, _ := ts.testRequest(t, "POST", "/updates/", reader)
	assert.Equal(t, http.StatusOK, statusCode)

	for _, req := range requests {
		path := fmt.Sprintf("/value/%s/%s", req.metricsType, req.metricsName)
		code, _ := ts.testRequest(t, "GET", path, nil)
		assert.Equal(t, http.StatusOK, code)
	}
}

func TestUpdateMetricsHandler(t *testing.T) {
	memStorage, err := repo.NewMetricsRepo()
	assert.NoError(t, err)

	tl := testLogger()
	ts := NewTestServer(memStorage, tl)

	gaugeMetrics := entity.Metrics{}

	statusCode, v := ts.testRequest(t, "POST", "/update/gauge/BuckHashSys/123.01", nil)
	err = json.Unmarshal(v, &gaugeMetrics)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, entity.Gauge(123.01), *gaugeMetrics.Value)

	counterMetrics := entity.Metrics{}

	statusCode, v = ts.testRequest(t, "POST", "/update/counter/PollCount/5", nil)
	err = json.Unmarshal(v, &counterMetrics)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, entity.Counter(5), *counterMetrics.Delta)

	statusCode, _ = ts.testRequest(t, "POST", "/update/superGauge/BuckHashSys/123.01", nil)
	assert.Equal(t, http.StatusNotImplemented, statusCode)

	statusCode, _ = ts.testRequest(t, "POST", "/update/counter/", nil)
	assert.Equal(t, http.StatusNotFound, statusCode)
}

func TestUpdateSpecificMetricsHandler(t *testing.T) {
	memStorage, err := repo.NewMetricsRepo()
	assert.NoError(t, err)

	tl := testLogger()
	ts := NewTestServer(memStorage, tl)

	var (
		tCounter entity.Counter = 5
		tGauge   entity.Gauge   = 123.01
	)

	type request struct {
		metricsType        string
		metricsName        string
		metricsValue       *entity.Gauge
		metricsDelta       *entity.Counter
		expectedStatusCode int
	}

	requests := []request{
		{
			metricsType:        Counter,
			metricsName:        "PollCount",
			metricsDelta:       &tCounter,
			expectedStatusCode: http.StatusOK,
		},
		{
			metricsType:        Gauge,
			metricsName:        "BuckHashSys",
			metricsValue:       &tGauge,
			expectedStatusCode: http.StatusOK,
		},
		//{
		//	metricsType:        "superGauge",
		//	metricsName:        "BuckHashSys",
		//	metricsValue:       &tGauge,
		//	expectedStatusCode: http.StatusNotImplemented,
		//},
		//{
		//	metricsType:        Counter,
		//	metricsName:        "",
		//	metricsDelta:       &tCounter,
		//	expectedStatusCode: http.StatusNotFound,
		//},
	}

	for _, req := range requests {
		var metricsVal float64
		if req.metricsDelta != nil {
			metricsVal = float64(*req.metricsDelta)
		} else if req.metricsValue != nil {
			metricsVal = float64(*req.metricsValue)
		}
		path := fmt.Sprintf("/update/%s/%s/%v", req.metricsType, req.metricsName, metricsVal)
		statusCode, b := ts.testRequest(t, "POST", path, nil)
		tempMetrics := entity.Metrics{}
		err = json.Unmarshal(b, &tempMetrics)
		assert.NoError(t, err)
		assert.Equal(t, req.expectedStatusCode, statusCode)
	}
}

func TestGetSomeMetricsHandler(t *testing.T) {

}

func TestGetSpecificMetricsHandler(t *testing.T) {

}

func TestPingHandler(t *testing.T) {
	memStorage, err := repo.NewMetricsRepo()
	assert.NoError(t, err)

	tl := testLogger()
	ts := NewTestServer(memStorage, tl)
	statusCode, _ := ts.testRequest(t, "GET", "/ping", nil)
	assert.Equal(t, http.StatusOK, statusCode)
}
