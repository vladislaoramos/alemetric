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

func (w *Worker) SendMetrics(ticker *time.Ticker) {
	for {
		<-ticker.C

		var wg sync.WaitGroup
		tasks := make(chan entity.Metrics)

		for i := 0; i < w.rateLimit; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				w.worker(tasks)
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

			task := entity.Metrics{
				ID:    name,
				MType: fieldType,
				Delta: valCounter,
				Value: valGauge,
			}

			time.Sleep(time.Millisecond * 500)

			tasks <- task
			w.l.Info(fmt.Sprintf("Metrics %s added to jobs list", name))
		}

		close(tasks)
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

func (w *Worker) worker(tasks chan entity.Metrics) {
	for task := range tasks {
		w.sendMetrics(task.ID, task.MType, task.Delta, task.Value)
	}
}
