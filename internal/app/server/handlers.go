package server

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/internal/usecase"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
	"net/http"
	"strings"
)

func getMetricsHandler(tool *usecase.ToolUseCase, l logger.LogInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		names, err := tool.GetMetricsNames()
		if err != nil {
			l.Error(fmt.Sprintf("Handlers - GetMetrics - Error: %s", err.Error()))
			errorHandler(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(strings.Join(names, "\n")))
	}
}

func updateSeveralMetricsHandler(tool *usecase.ToolUseCase, l logger.LogInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var items []entity.Metrics
		if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
			http.Error(w, "error decoding several metrics", http.StatusBadRequest)
			return
		}

		for _, item := range items {
			if err := tool.StoreMetrics(item); err != nil {
				l.Error(fmt.Sprintf("Handlers - UpdateSeveralMetrics - Error: %s", err.Error()))
				errorHandler(w, err)
				return
			}
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}

func updateMetricsHandler(tool *usecase.ToolUseCase, l logger.LogInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics entity.Metrics
		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			http.Error(w, "error decoding metrics", http.StatusBadRequest)
			return
		}

		if err := tool.StoreMetrics(metrics); err != nil {
			l.Error(fmt.Sprintf("Handlers - UpdateMetrics - Error: %s", err.Error()))
			errorHandler(w, err)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}

func updateSpecificMetricsHandler(tool *usecase.ToolUseCase, l logger.LogInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricsType := chi.URLParam(r, "metricsType")
		metricsName := chi.URLParam(r, "metricsName")
		metricsValue := chi.URLParam(r, "metricsValue")

		switch metricsType {
		case Gauge:
			value, err := entity.ParseGaugeMetrics(metricsValue)
			if err != nil {
				l.Error(fmt.Sprintf("Handlers - UpdateSpecificMetrics - ParseGauge: %s", err.Error()))
				http.Error(w, "error parsing gauge metric during update", http.StatusBadRequest)
			}

			metrics := entity.Metrics{
				ID:    metricsName,
				MType: value.Type(),
				Value: &value,
			}

			err = tool.StoreMetrics(metrics)
			if err != nil {
				l.Error(fmt.Sprintf("Handlers - UpdateSpecificMetrics - Error: %s", err.Error()))
				errorHandler(w, err)
				return
			}
		case Counter:
			value, err := entity.ParseCounterMetrics(metricsValue)
			if err != nil {
				l.Error(fmt.Sprintf("Handlers - UpdateSpecificMetrics - ParseCounter: %s", err.Error()))
				http.Error(w, "error parsing counter metric during update", http.StatusBadRequest)
			}

			metrics := entity.Metrics{
				ID:    metricsName,
				MType: value.Type(),
				Delta: &value,
			}

			err = tool.StoreMetrics(metrics)
			if err != nil {
				l.Error(fmt.Sprintf("Handlers - UpdateSpecificMetrics - Error: %s", err.Error()))
				errorHandler(w, err)
				return
			}
		default:
			l.Error(fmt.Sprintf("Handlers - UpdateSpecificMetrics - Metrics Type: %s", metricsType))
			http.Error(w, "metrics type not found", http.StatusNotImplemented)
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
	}
}

func getSomeMetricsHandler(tool *usecase.ToolUseCase, l logger.LogInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics entity.Metrics
		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			http.Error(w, "error decoding metrics during get", http.StatusBadRequest)
			return
		}

		value, err := tool.GetMetrics(metrics)
		if err != nil {
			l.Error(fmt.Sprintf("Handlers - GetSomeMetrics - Error: %s", err.Error()))
			errorHandler(w, err)
			return
		}

		resp, err := json.Marshal(value)
		if err != nil {
			l.Error(fmt.Sprintf("Handlers - GetSomeMetrics - Response Marshaling: %s", err.Error()))
			errorHandler(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func getSpecificMetricsHandler(tool *usecase.ToolUseCase, l logger.LogInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricsType := chi.URLParam(r, "metricsType")
		metricsName := chi.URLParam(r, "metricsName")

		metrics := entity.Metrics{
			ID:    metricsName,
			MType: metricsType,
		}

		res, err := tool.GetMetrics(metrics)
		if err != nil {
			l.Error(fmt.Sprintf("Handlers - GetSpecificMetrics - Error: %s", err.Error()))
			http.Error(w, "metrics not found", http.StatusNotFound)
			return
		}

		switch metricsType {
		case Gauge:
			w.Write([]byte(fmt.Sprintf("%g", *res.Value)))
		case Counter:
			w.Write([]byte(fmt.Sprintf("%d", *res.Delta)))
		default:
			http.Error(w, "metrics type not found", http.StatusNotImplemented)
		}
	}
}

func pingHandler(tool *usecase.ToolUseCase, l logger.LogInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tool.PingRepo(r.Context()); err != nil {
			l.Error(fmt.Sprintf("Handlers - PignHandlers - DB Connection Error: %s", err.Error()))
			http.Error(w, "error db connection", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
