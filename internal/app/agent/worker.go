package agent

import (
	"fmt"
	"github.com/vladislaoramos/alemetric/internal/usecase"
	"reflect"
	"strings"
	"sync"
	"time"

	logger "github.com/vladislaoramos/alemetric/pkg/log"

	"github.com/vladislaoramos/alemetric/internal/entity"
)

type Worker struct {
	webAPI           WebAPIAgent
	metrics          *Metrics
	metricsNames     []string
	l                logger.LogInterface
	rateLimitCounter uint
}

func NewWorker(
	l logger.LogInterface,
	metrics *Metrics,
	metricsNames []string,
	webAPI WebAPIAgent,
	limit uint) *Worker {
	return &Worker{
		l:                l,
		metrics:          metrics,
		metricsNames:     metricsNames,
		webAPI:           webAPI,
		rateLimitCounter: limit,
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

		var workersNum int
		if w.rateLimitCounter > 0 {
			workersNum = int(w.rateLimitCounter)
		} else {
			w.l.Fatal(
				fmt.Sprintf(
					"The current number of workers is %d. It must be positive and greater than 0",
					w.rateLimitCounter))
		}

		for i := 0; i < workersNum; i++ {
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
			case usecase.Counter:
				val := entity.Counter(field.Int())
				valCounter = &val
			case usecase.Gauge:
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

			tasks <- task
			w.l.Info(fmt.Sprintf("Metrics %s added to jobs list", name))
		}

		close(tasks)
		wg.Wait()
	}
}

func (w *Worker) sendMetrics(name, mType string, counter *entity.Counter, gauge *entity.Gauge) {
	w.l.Info(fmt.Sprintf("Metrics %s is sending", name))

	var c entity.Counter
	if counter != nil {
		c = *counter
	}

	var g entity.Gauge
	if gauge != nil {
		g = *gauge
	}

	err := w.webAPI.SendMetrics(name, mType, counter, gauge)
	if err != nil {
		w.l.Error(
			fmt.Sprintf(
				"error sending metrics conflict: %v; metricName: %s metricType: %s delta: %d value: %v",
				err, name, mType, c, g))
	}
}

func (w *Worker) worker(tasks chan entity.Metrics) {
	for task := range tasks {
		w.sendMetrics(task.ID, task.MType, task.Delta, task.Value)
	}
}
