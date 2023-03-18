package agent

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	logger "github.com/vladislaoramos/alemetric/pkg/log"

	"github.com/vladislaoramos/alemetric/internal/entity"
)

type Worker struct {
	webAPI       WebAPIAgent
	metrics      *Metrics
	metricsNames []string
	l            logger.LogInterface
	rateLimit    int
}

func NewWorker(
	l logger.LogInterface,
	metrics *Metrics,
	metricsNames []string,
	webAPI WebAPIAgent,
	limit int) *Worker {
	return &Worker{
		l:            l,
		metrics:      metrics,
		metricsNames: metricsNames,
		webAPI:       webAPI,
		rateLimit:    limit,
	}
}

func (w *Worker) UpdateMetrics(ticker *time.Ticker) {
	for {
		<-ticker.C
		w.metrics.CollectMetrics()
		w.metrics.CollectAdditionalMetrics()
		w.l.Info("Metrics updated")
	}
}

func (w *Worker) SendMetrics(ticker *time.Ticker) {
	for {
		<-ticker.C

		jobCh := make(chan entity.Metrics, 1)
		for i := 0; i < w.rateLimit; i++ {
			go func() {
				for job := range jobCh {
					w.l.Info(fmt.Sprintf("Job with metrics %s is extracting from job channel", job.ID))
					go w.sendMetrics(job.ID, job.MType, job.Delta, job.Value)
				}
			}()
		}

		var jobs []entity.Metrics

		for _, name := range w.metricsNames {
			field := reflect.Indirect(reflect.ValueOf(w.metrics)).FieldByName(name)
			if !field.IsValid() {
				w.l.Error(fmt.Sprintf("Field `%s` is not valid", name))
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
				w.l.Error(fmt.Sprintf("Type of the metrics field `%s` is invalid", fieldType))
				continue
			}

			time.Sleep(time.Second)
			//go func(metricsName, metricsType string, delta *entity.Counter, value *entity.Gauge) {
			//	err := w.webAPI.SendMetrics(metricsName, metricsType, delta, value)
			//	if err != nil {
			//		w.l.Error(fmt.Sprintf("error sending metrics conflict: %v; metricName: %s metricType: %s delta: %v value: %v", err,
			//			metricsName, metricsType, delta, value))
			//	}
			//}(name, fieldType, valCounter, valGauge)

			jobs = append(jobs, entity.Metrics{
				ID:    name,
				MType: fieldType,
				Delta: valCounter,
				Value: valGauge,
			})
			w.l.Info(fmt.Sprintf("Metrics %s added to jobs list", name))
		}

		for _, j := range jobs {
			w.l.Info(fmt.Sprintf("Job with metrics %s is adding to job channel", j.ID))
			jobCh <- j
		}
	}
}

func (w *Worker) sendMetrics(name, mType string, counter *entity.Counter, gauge *entity.Gauge) {
	go func(metricsName, metricsType string, delta *entity.Counter, value *entity.Gauge) {
		w.l.Info(fmt.Sprintf("Metrics %s is sending", name))
		err := w.webAPI.SendMetrics(metricsName, metricsType, delta, value)
		if err != nil {
			w.l.Error(fmt.Sprintf("error sending metrics conflict: %v; metricName: %s metricType: %s delta: %v value: %v", err,
				metricsName, metricsType, delta, value))
		}
	}(name, mType, counter, gauge)
}
