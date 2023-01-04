package agent

import (
	"github.com/vladislaoramos/alemetric/pkg/log"
	"time"
)

type Worker struct {
	l       logger.LogInterface
	metrics *Metrics
	webAPI  *Service
}

func NewWorker(l logger.LogInterface, metrics *Metrics, webAPI *Service) *Worker {
	return &Worker{
		l:       l,
		metrics: metrics,
		webAPI:  webAPI,
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

		err := w.webAPI.SendGaugeMetric("Alloc", w.metrics.storage.Alloc)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("BuckHashSys", w.metrics.storage.BuckHashSys)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("Frees", w.metrics.storage.Frees)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("GCCPUFraction", w.metrics.storage.GCCPUFraction)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("GCSys", w.metrics.storage.GCSys)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("HeapAlloc", w.metrics.storage.HeapAlloc)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("HeapIdle", w.metrics.storage.HeapIdle)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("HeapInuse", w.metrics.storage.HeapInuse)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("HeapObjects", w.metrics.storage.HeapObjects)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("HeapReleased", w.metrics.storage.HeapReleased)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("HeapSys", w.metrics.storage.HeapSys)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("LastGC", w.metrics.storage.LastGC)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("Lookups", w.metrics.storage.Lookups)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("MCacheInuse", w.metrics.storage.MCacheInuse)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("MCacheSys", w.metrics.storage.MCacheSys)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("MSpanInuse", w.metrics.storage.MSpanInuse)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("MSpanSys", w.metrics.storage.MSpanSys)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("Mallocs", w.metrics.storage.Mallocs)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("NextGC", w.metrics.storage.NextGC)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("NumForcedGC", w.metrics.storage.NumForcedGC)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("NumGC", w.metrics.storage.NumGC)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("OtherSys", w.metrics.storage.OtherSys)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("PauseTotalNs", w.metrics.storage.PauseTotalNs)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("StackInuse", w.metrics.storage.StackInuse)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("StackSys", w.metrics.storage.StackSys)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("Sys", w.metrics.storage.Sys)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("TotalAlloc", w.metrics.storage.TotalAlloc)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendGaugeMetric("RandomValue", w.metrics.RandomValue)
		if err != nil {
			w.l.Error(err.Error())
		}

		err = w.webAPI.SendCounterMetric("PollCount", w.metrics.PollCount)
		if err != nil {
			w.l.Error(err.Error())
		}
	}
}
