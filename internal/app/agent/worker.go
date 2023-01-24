package agent

import (
	"fmt"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/pkg/log"
	"reflect"
	"strings"
	"time"
)

type Worker struct {
	webAPI       WebAPIAgent
	metrics      *Metrics
	metricsNames []string
	l            logger.LogInterface
}

func NewWorker(
	l logger.LogInterface,
	metrics *Metrics,
	metricsNames []string,
	webAPI WebAPIAgent) *Worker {
	return &Worker{
		l:            l,
		metrics:      metrics,
		metricsNames: metricsNames,
		webAPI:       webAPI,
	}
}

func (w *Worker) UpdateMetrics(ticker *time.Ticker) {
	for {
		<-ticker.C
		w.metrics.CollectMetrics()
		w.l.Info("metrics updated")
	}
}

func (w *Worker) SendMetrics(ticker *time.Ticker) {
	for {
		<-ticker.C

		for _, name := range w.metricsNames {
			field := reflect.Indirect(reflect.ValueOf(w.metrics)).FieldByName(name)
			if !field.IsValid() {
				w.l.Error(fmt.Sprintf("field `%s` is not valid", name))
				continue
			}

			fieldType := strings.ToLower(field.Type().Name())

			var (
				valCounter *entity.Counter
				valGauge   *entity.Gauge
			)

			switch fieldType {
			case "counter":
				val := entity.Counter(field.Int())
				valCounter = &val
			case "gauge":
				val := entity.Gauge(field.Float())
				valGauge = &val
			default:
				w.l.Error(fmt.Sprintf("type of the metrics field %s is invalid", fieldType))
				continue
			}

			go func(metricsName, metricsType string, delta *entity.Counter, value *entity.Gauge) {
				err := w.webAPI.SendMetrics(metricsName, metricsType, delta, value)
				if err != nil {
					w.l.Error(fmt.Sprintf("error sending metrics conflict: %s", err))
				}
			}(name, fieldType, valCounter, valGauge)
		}
	}
}
