package server

import (
	"os"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/vladislaoramos/alemetric/internal/usecase"
	"github.com/vladislaoramos/alemetric/internal/usecase/mocks"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

func testLogger() *logger.Logger {
	f, err := os.OpenFile("/tmp/test_log_server", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0777)
	l := logger.New("debug", f)
	if err != nil {
		l.Fatal("unable to open file for log")
	}

	return l
}

func metricsTool(t *testing.T) (*usecase.ToolUseCase, *mocks.MetricsRepo) {
	log := testLogger()
	mockCtl := gomock.NewController(t)
	defer mockCtl.Finish()

	mockRepo := mocks.NewMetricsRepo(t)

	mockTool := usecase.NewMetricsTool(
		mockRepo,
		log,
	)

	return mockTool, mockRepo
}

func TestGetMetricsHandler(t *testing.T) {

}

func TestUpdateSeveralMetricsHandler(t *testing.T) {

}

func TestUpdateMetricsHandler(t *testing.T) {

}

func TestUpdateSpecificMetricsHandler(t *testing.T) {

}

func TestGetSomeMetricsHandler(t *testing.T) {

}

func TestGetSpecificMetricsHandler(t *testing.T) {

}
