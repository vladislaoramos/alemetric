package agent

import (
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
	"github.com/vladislaoramos/alemetric/internal/entity"
)

// WebAPIClient implements the client web-application for Agent.
type WebAPIClient struct {
	client *resty.Client
	Key    string
}

func NewWebAPI(client *resty.Client, key string) *WebAPIClient {
	return &WebAPIClient{
		client: client,
		Key:    key,
	}
}

// SendMetrics sends a client request for a metrics update to the server.
func (wc *WebAPIClient) SendMetrics(
	metricsName,
	metricsType string,
	delta *entity.Counter,
	value *entity.Gauge,
) error {
	body := entity.Metrics{
		ID:    metricsName,
		MType: metricsType,
		Delta: delta,
		Value: value,
	}

	body.SignData("agent", wc.Key)

	// respBody, _ := json.Marshal(body)
	// log.Printf("req body: %s", string(respBody))

	resp, err := wc.client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		Post("/update/")
	if err != nil {
		return fmt.Errorf("cannot send metrics from agent: %w", err)
	}

	status := resp.StatusCode()
	if status != http.StatusOK {
		return fmt.Errorf("sending metrics from agent with not successful status code: %d", status)
	}

	return nil
}

// SendSeveralMetrics sends a client request for several metrics update to the server.
func (wc *WebAPIClient) SendSeveralMetrics(items []entity.Metrics) error {
	resp, err := wc.client.
		R().
		SetHeader("Content-Type", "application/json").
		SetBody(items).
		Post("/updates/")
	if err != nil {
		return fmt.Errorf("cannot send several metrics from agent: %w", err)
	}

	status := resp.StatusCode()
	if status != http.StatusOK {
		return fmt.Errorf("sending several metrics from agent with not successful status code: %d", status)
	}
	return nil
}
