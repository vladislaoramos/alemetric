package agent

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
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
		w.l.Info("Metrics updated")
	}
}

func (w *Worker) UpdateAdditionalMetrics(ticker *time.Ticker) {
	for {
		<-ticker.C
		w.metrics.CollectAdditionalMetrics()
		w.l.Info("Additional metrics updated")
	}
}

func (w *Worker) SendMetrics_(ticker *time.Ticker) {
	// start := time.Now()
	for {
		<-ticker.C

		wg := sync.WaitGroup{}
		jobCh := make(chan entity.Metrics)

		for i := 0; i < w.rateLimit; i++ {
			go func() {
				for job := range jobCh {
					w.l.Info(fmt.Sprintf("Job with metrics %s is extracting from job channel", job.ID))
					// w.l.Info(fmt.Sprintf("Time after start: %v", time.Now().Sub(start)))
					//wg.Done()
					w.sendMetrics(job.ID, job.MType, job.Delta, job.Value)
					wg.Done()
				}
			}()
		}

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

			job := entity.Metrics{
				ID:    name,
				MType: fieldType,
				Delta: valCounter,
				Value: valGauge,
			}

			//wg.Add(1)
			jobCh <- job
			wg.Add(1)
			w.l.Info(fmt.Sprintf("Metrics %s added to jobs list", name))
		}
		wg.Wait()
	}
}

func (w *Worker) sendMetrics(name, mType string, counter *entity.Counter, gauge *entity.Gauge) {
	w.l.Info(fmt.Sprintf("Metrics %s is sending", name))
	err := w.webAPI.SendMetrics(name, mType, counter, gauge)
	if err != nil {
		w.l.Error(
			fmt.Sprintf(
				"error sending metrics conflict: %v; metricName: %s metricType: %s delta: %v value: %v",
				err, name, mType, counter, gauge))
	}
}

func (w *Worker) SendMetrics(ticker *time.Ticker) {
	for {
		<-ticker.C

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
			go func(metricsName, metricsType string, delta *entity.Counter, value *entity.Gauge) {
				err := w.webAPI.SendMetrics(metricsName, metricsType, delta, value)
				if err != nil {
					w.l.Error(fmt.Sprintf("error sending metrics conflict: %v; metricName: %s metricType: %s delta: %v value: %v", err,
						metricsName, metricsType, delta, value))
				}
			}(name, fieldType, valCounter, valGauge)
		}
	}
}
