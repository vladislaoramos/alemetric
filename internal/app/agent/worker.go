package agent

import (
	"fmt"
	"github.com/vladislaoramos/alemetric/pkg/log"
	"reflect"
	"strings"
	"time"
)

type Worker struct {
	l            logger.LogInterface
	metrics      *Metrics
	metricsNames []string
	webAPI       *Service
}

func NewWorker(
	l logger.LogInterface,
	metrics *Metrics,
	metricsNames []string,
	webAPI *Service) *Worker {
	return &Worker{
		l:            l,
		metrics:      metrics,
		webAPI:       webAPI,
		metricsNames: metricsNames,
	}
}

func (w *Worker) UpdateMetrics(ticker *time.Ticker) {
	for {
		<-ticker.C
		w.metrics.UpdateMetrics()
		w.l.Info("metrics updated")
	}
}

func (w *Worker) SendMetrics(ticker *time.Ticker) {
	for {
		<-ticker.C

		for _, mName := range w.metricsNames {
			r := reflect.ValueOf(w.metrics)
			f := reflect.Indirect(r).FieldByName(mName)
			if !f.IsValid() {
				w.l.Error(fmt.Sprintf("field `%s` is not valid", mName))
				continue
			}
			go func(metricName, metricType string, metricValue interface{}) {
				err := w.webAPI.SendMetrics(metricName, metricType, metricValue)
				if err != nil {
					w.l.Error(err.Error())
				}
			}(mName, strings.ToLower(f.Type().Name()), f)
		}
	}
}
