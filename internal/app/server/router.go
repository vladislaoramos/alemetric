package server

import (
	"fmt"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"io"
	"net/http"
	"strings"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

const (
	Counter = "counter"
	Gauge   = "gauge"
)

func NewRouter(handler *chi.Mux, repo MetricsRepo) {
	handler.Use(middleware.RealIP)
	handler.Use(middleware.Logger)
	handler.Use(middleware.Recoverer)
	handler.Use(middleware.RequestID)

	handler.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, strings.Join(repo.GetMetricsNames(), "\n"))
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		}
	})

	handler.Route("/update", func(r chi.Router) {
		r.Post(
			"/{metricsType}/{metricsName}/{metricsValue}",
			func(w http.ResponseWriter, r *http.Request) {
				metricsType := chi.URLParam(r, "metricsType")
				metricsName := chi.URLParam(r, "metricsName")
				metricsValue := chi.URLParam(r, "metricsValue")

				switch metricsType {
				case Counter:
					value, err := entity.ParseCounterMetrics(metricsValue)
					if err != nil {
						http.Error(w, "bad request", http.StatusBadRequest)
						return
					}
					err = repo.StoreCounterMetrics(metricsName, value)
					if err != nil {
						http.Error(w, "some problem with storage", http.StatusInternalServerError)
					}
				case Gauge:
					value, err := entity.ParseGaugeMetrics(metricsValue)
					if err != nil {
						http.Error(w, "bad request", http.StatusBadRequest)
						return
					}
					err = repo.StoreGaugeMetrics(metricsName, value)
					if err != nil {
						http.Error(w, "some problem with storage", http.StatusInternalServerError)
					}
				default:
					http.Error(w, "metrics type is not found", http.StatusNotImplemented)
				}
				w.WriteHeader(http.StatusOK)
			},
		)
	})

	handler.Route("/value", func(r chi.Router) {
		r.Get("/{metricType}/{metricName}", func(w http.ResponseWriter, r *http.Request) {
			metricType := chi.URLParam(r, "metricType")
			metricName := chi.URLParam(r, "metricName")

			switch metricType {
			case Counter, Gauge:
				value, err := repo.GetMetrics(metricName)
				if err != nil {
					http.Error(w, "metrics is not found", http.StatusNotFound)
					return
				}
				w.Write([]byte(fmt.Sprintf("%d", value)))
			default:
				http.Error(w, "metrics type is not found", http.StatusNotImplemented)
			}
		})
	})
}
