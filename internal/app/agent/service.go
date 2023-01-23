package agent

import (
	"fmt"
	"github.com/go-resty/resty/v2"
	"net/http"
)

type WebAPIClient struct {
	client *resty.Client
}

func NewWebAPI(client *resty.Client) *WebAPIClient {
	return &WebAPIClient{
		client: client,
	}
}

func (wc *WebAPIClient) SendMetrics(metricsName, metricsType string, metricsValue interface{}) error {
	resp, err := wc.client.R().SetHeader("Content-Type", "text/plain").
		Post(fmt.Sprintf("/update/%s/%s/%v", metricsType, metricsName, metricsValue))
	if err != nil {
		return fmt.Errorf("cannot send metrics because of error: %w", err)
	}

	status := resp.StatusCode()
	if status != http.StatusOK {
		return fmt.Errorf("error sending metrics with status code: %d", status)
	}

	return nil
}
