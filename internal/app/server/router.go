package server

import (
	"encoding/json"
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
			"/",
			func(w http.ResponseWriter, r *http.Request) {
				var metrics entity.Metrics
				if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
					http.Error(w, "error with updating metrics", http.StatusBadRequest)
					return
				}

				switch metrics.MType {
				case Counter:
					if err := repo.StoreCounterMetrics(metrics.ID, *metrics.Delta); err != nil {
						http.Error(w, "some problem with storage", http.StatusInternalServerError)
					}
				case Gauge:
					if err := repo.StoreGaugeMetrics(metrics.ID, *metrics.Value); err != nil {
						http.Error(w, "some problem with storage", http.StatusInternalServerError)
					}
				}

				w.WriteHeader(http.StatusOK)
			},
		)
		r.Post(
			"/{metricsType}/{metricsName}/{metricsValue}",
			func(w http.ResponseWriter, r *http.Request) {
				metricsType := chi.URLParam(r, "metricsType")
				metricsName := chi.URLParam(r, "metricsName")
				metricsValue := chi.URLParam(r, "metricsValue")

				switch metricsType {
				case Gauge:
					value, err := entity.ParseGaugeMetrics(metricsValue)
					if err != nil {
						http.Error(w, "bad value type", http.StatusBadRequest)
						return
					}
					err = repo.StoreGaugeMetrics(metricsName, value)
					if err != nil {
						http.Error(w, "some problem with storage", http.StatusInternalServerError)
					}
				case Counter:
					value, err := entity.ParseCounterMetrics(metricsValue)
					if err != nil {
						http.Error(w, "bad value type", http.StatusBadRequest)
						return
					}
					err = repo.StoreCounterMetrics(metricsName, value)
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
		r.Post(
			"/",
			func(w http.ResponseWriter, r *http.Request) {
				var metrics entity.Metrics
				if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
					http.Error(w, "error with getting metrics", http.StatusBadRequest)
					return
				}

				value, err := repo.GetMetrics(metrics.ID)
				if err != nil {
					http.Error(w, "metrics not found", http.StatusNotFound)
					return
				}

				switch metrics.MType {
				case Counter:
					curVal := value.(entity.Counter)
					metrics.Delta = &curVal
				case Gauge:
					curVal := value.(entity.Gauge)
					metrics.Value = &curVal
				}

				resp, err := json.Marshal(metrics)
				if err != nil {
					http.Error(w, "server error", http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write(resp)
			},
		)
		r.Get("/{metricsType}/{metricsName}", func(w http.ResponseWriter, r *http.Request) {
			metricsType := chi.URLParam(r, "metricsType")
			metricsName := chi.URLParam(r, "metricsName")

			switch metricsType {
			case Gauge:
				value, err := repo.GetMetrics(metricsName)
				if err != nil {
					http.Error(w, "metrics is not found", http.StatusNotFound)
					return
				}
				w.Write([]byte(fmt.Sprintf("%g", value)))
			case Counter:
				value, err := repo.GetMetrics(metricsName)
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
