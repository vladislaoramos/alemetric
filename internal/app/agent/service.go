package agent

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type Service struct {
	client *resty.Client
}

func NewAPI(client *resty.Client) *Service {
	return &Service{
		client: client,
	}
}

func (s *Service) SendMetrics(metricsName, metricsType string, metricsValue interface{}) error {
	resp, err := s.client.
		R().
		SetHeader("Content-Type", "text/plain").
		Post(
			fmt.Sprintf("/update/%s/%s/%v", metricsType, metricsName, metricsValue),
		)
	if err != nil {
		return fmt.Errorf("cant't send metrics: %w", err)
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("cant't send metrics. Status code <> 200")
	}
	return nil
}
