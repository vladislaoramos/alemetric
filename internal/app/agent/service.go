package agent

import (
	"fmt"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/pkg/log"
	"net/http"
)

type Service struct {
	l      logger.LogInterface
	host   string
	port   string
	client *http.Client
}

func NewAPI(l logger.LogInterface, host, port string) *Service {
	return &Service{
		l:    l,
		host: host,
		port: port,
		client: &http.Client{
			Transport: &http.Transport{},
		},
	}
}

func makeAddr(host, port string) string {
	return "http://" + host + ":" + port + "/"
}

func (s *Service) SendGaugeMetric(metricName string, metricValue entity.Gauge) error {
	addr := makeAddr(s.host, s.port)

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%supdate/gauge/%s/%f", addr, metricName, metricValue),
		nil,
	)
	if err != nil {
		return fmt.Errorf("SendGaugeMetric with error: %w", err)
	}

	request.Header.Add("Content-Type", "text/plain")
	response, err := s.client.Do(request)
	if err != nil {
		return fmt.Errorf("SendGaugeMetric with error: %w", err)
	}

	defer response.Body.Close()
	s.l.Info(fmt.Sprintf("send with status: %s", response.Status))

	return nil
}

func (s *Service) SendCounterMetric(metricName string, metricValue entity.Counter) error {
	addr := makeAddr(s.host, s.port)

	request, err := http.NewRequest(
		http.MethodPost,
		fmt.Sprintf("%supdate/counter/%s/%d", addr, metricName, metricValue),
		nil,
	)
	if err != nil {
		return fmt.Errorf("SendCounterMetric with error: %w", err)
	}

	request.Header.Add("Content-Type", "text/plain")

	response, err := s.client.Do(request)
	if err != nil {
		return fmt.Errorf("SendCounterMetric with error: %w", err)
	}

	defer response.Body.Close()

	s.l.Info(fmt.Sprintf("send with status: %s", response.Status))

	return nil
}
