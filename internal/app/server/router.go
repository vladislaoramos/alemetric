package server

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"io"
	"net/http"
	"strings"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

func NewRouter(handler *chi.Mux, repo MetricsRepo) {
	handler.Use(middleware.RequestID)
	handler.Use(middleware.RealIP)
	handler.Use(middleware.Logger)
	handler.Use(middleware.Recoverer)

	handler.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := io.WriteString(w, strings.Join(repo.GetMetricsNames(), "\n"))
		if err != nil {
			panic(err)
		}
	})

	// update
	handler.Route("/update", func(r chi.Router) {
		r.Post(
			"/{metricType}/{metricName}/{metricValue}",
			func(w http.ResponseWriter, r *http.Request) {
				metricType := chi.URLParam(r, "metricType")
				metricName := chi.URLParam(r, "metricName")
				metricValue := chi.URLParam(r, "metricValue")

				switch metricType {
				case Gauge:
					value, err := entity.ParseGaugeMetrics(metricValue)
					if err != nil {
						http.Error(w, "bad value type", http.StatusBadRequest)
						return
					}
					err = repo.StoreGaugeMetrics(metricName, value)
					if err != nil {
						http.Error(w, "some problem with storage", http.StatusInternalServerError)
					}
				case Counter:
					value, err := entity.ParseCounterMetrics(metricValue)
					if err != nil {
						http.Error(w, "bad value type", http.StatusBadRequest)
						return
					}
					err = repo.StoreCounterMetrics(metricName, value)
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

	// value
	handler.Route("/value", func(r chi.Router) {
		r.Get("/{metricType}/{metricName}", func(w http.ResponseWriter, r *http.Request) {
			metricType := chi.URLParam(r, "metricType")
			metricName := chi.URLParam(r, "metricName")

			switch metricType {
			case Gauge:
				value, err := repo.GetMetrics(metricName)
				if err != nil {
					http.Error(w, "metrics is not found", http.StatusNotFound)
					return
				}
				w.Write([]byte(fmt.Sprintf("%g", value)))
			case Counter:
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
