package agent

import (
	"fmt"
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
			//r := reflect.ValueOf(w.metrics)
			value := reflect.Indirect(reflect.ValueOf(w.metrics)).FieldByName(name)
			if !value.IsValid() {
				w.l.Error(fmt.Sprintf("field `%s` is not valid", name))
				continue
			}

			go func(metricsName, metricsType string, metricsValue interface{}) {
				err := w.webAPI.SendMetrics(metricsName, metricsType, metricsValue)
				if err != nil {
					w.l.Error(err.Error())
				}
			}(name, strings.ToLower(value.Type().Name()), value)
		}
	}
}
