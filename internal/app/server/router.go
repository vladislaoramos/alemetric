package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/internal/usecase"
	"github.com/vladislaoramos/alemetric/pkg/log"
	"io"
	"net/http"
	"strings"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

func NewRouter(handler *chi.Mux, tool *usecase.ToolUseCase, l logger.LogInterface) {
	handler.Use(middleware.RequestID)
	handler.Use(middleware.RealIP)
	handler.Use(middleware.Logger)
	handler.Use(middleware.Recoverer)

	handler.Get("/", func(w http.ResponseWriter, r *http.Request) {
		names, err := tool.GetMetricsNames()
		if err != nil {
			errorHandler(w, err)
			return
		}
		_, err = io.WriteString(w, strings.Join(names, "\n"))
		if err != nil {
			errorHandler(w, err)
			return
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

				if err := tool.StoreMetrics(metrics); err != nil {
					l.Error(err.Error())
					errorHandler(w, err)
					return
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
						l.Error(err.Error())
						http.Error(w, "parsing error", http.StatusBadRequest)
					}

					metrics := entity.Metrics{
						ID:    metricsName,
						MType: value.Type(),
						Value: &value,
					}

					err = tool.StoreMetrics(metrics)
					if err != nil {
						l.Error(err.Error())
						errorHandler(w, err)
						return
					}
				case Counter:
					value, err := entity.ParseCounterMetrics(metricsValue)
					if err != nil {
						l.Error(err.Error())
						http.Error(w, "parsing error", http.StatusBadRequest)
					}

					metrics := entity.Metrics{
						ID:    metricsName,
						MType: value.Type(),
						Delta: &value,
					}

					err = tool.StoreMetrics(metrics)
					if err != nil {
						l.Error(err.Error())
						errorHandler(w, err)
						return
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

				value, err := tool.GetMetrics(metrics)
				if err != nil {
					l.Error(err.Error())
					errorHandler(w, err)
					return
				}

				resp, err := json.Marshal(value)
				if err != nil {
					l.Error(err.Error())
					errorHandler(w, err)
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.Write(resp)
			},
		)
		r.Get("/{metricsType}/{metricsName}", func(w http.ResponseWriter, r *http.Request) {
			metricsType := chi.URLParam(r, "metricsType")
			metricsName := chi.URLParam(r, "metricsName")

			metrics := entity.Metrics{
				ID:    metricsName,
				MType: metricsType,
			}

			res, err := tool.GetMetrics(metrics)
			if err != nil {
				l.Error(err.Error())
				http.Error(w, "metrics is not found", http.StatusNotFound)
				return
			}

			switch metricsType {
			case Gauge:
				w.Write([]byte(fmt.Sprintf("%g", *res.Value)))
			case Counter:
				w.Write([]byte(fmt.Sprintf("%d", *res.Delta)))
			default:
				http.Error(w, "metrics type is not found", http.StatusNotImplemented)
			}
		})
	})
}
