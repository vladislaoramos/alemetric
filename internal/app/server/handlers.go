package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/internal/usecase"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

func getMetricsHandler(tool *usecase.ToolUseCase, l logger.LogInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		names, err := tool.GetMetricsNames(r.Context())
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
			http.Error(w, "error decoding several metrics: "+err.Error(), http.StatusBadRequest)
			return
		}

		for _, item := range items {
			if err := tool.StoreMetrics(r.Context(), item); err != nil {
				l.Error(fmt.Errorf("error with updating metrics: %w", err).Error())
				errorHandler(w, err)
				return
			}
		}

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

		if err := tool.StoreMetrics(r.Context(), metrics); err != nil {
			l.Error(fmt.Errorf("error with updating metrics: %w", err).Error())
			errorHandler(w, err)
			return
		}

		value, err := tool.GetMetrics(r.Context(), metrics)
		if err != nil {
			l.Error(fmt.Errorf("error with getting updated metrics: %w", err).Error())
			errorHandler(w, err)
			return
		}

		resp, err := json.Marshal(value)
		if err != nil {
			l.Error(fmt.Errorf("error with marshalling metrics before sending the response: %w", err).Error())
			errorHandler(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func updateSpecificMetricsHandler(tool *usecase.ToolUseCase, l logger.LogInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		metricsType := chi.URLParam(r, "metricsType")
		metricsName := chi.URLParam(r, "metricsName")
		metricsValue := chi.URLParam(r, "metricsValue")

		var metrics entity.Metrics
		switch metricsType {
		case Gauge:
			value, err := entity.ParseGaugeMetrics(metricsValue)
			if err != nil {
				l.Error(err.Error())
				http.Error(w, "parsing error", http.StatusBadRequest)
				return
			}

			metrics = entity.Metrics{
				ID:    metricsName,
				MType: value.Type(),
				Value: &value,
			}

			err = tool.StoreMetrics(r.Context(), metrics)
			if err != nil {
				l.Error(fmt.Sprintf("Handlers - UpdateSpecificMetrics - Error: %s", err.Error()))
				errorHandler(w, err)
				return
			}
		case Counter:
			value, err := entity.ParseCounterMetrics(metricsValue)
			if err != nil {
				l.Error(err.Error())
				http.Error(w, "parsing error", http.StatusBadRequest)
				return
			}

			metrics = entity.Metrics{
				ID:    metricsName,
				MType: value.Type(),
				Delta: &value,
			}

			err = tool.StoreMetrics(r.Context(), metrics)
			if err != nil {
				l.Error(fmt.Sprintf("Handlers - UpdateSpecificMetrics - Error: %s", err.Error()))
				errorHandler(w, err)
				return
			}
		default:
			l.Error(fmt.Sprintf("Handlers - UpdateSpecificMetrics - Metrics Type: %s", metricsType))
			http.Error(w, "metrics type not found", http.StatusNotImplemented)
		}

		resp, err := json.Marshal(metrics)
		if err != nil {
			l.Error(err.Error())
			errorHandler(w, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
	}
}

func getSomeMetricsHandler(tool *usecase.ToolUseCase, l logger.LogInterface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var metrics entity.Metrics
		if err := json.NewDecoder(r.Body).Decode(&metrics); err != nil {
			http.Error(w, "error decoding metrics during get", http.StatusBadRequest)
			return
		}

		value, err := tool.GetMetrics(r.Context(), metrics)
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

		fmt.Printf(`getSomeMetrics:
			original struct: %+v
			json: %s
		`, value, string(resp))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
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

		res, err := tool.GetMetrics(r.Context(), metrics)
		if err != nil {
			l.Error(err.Error())
			errorHandler(w, err)
			return
		}

		var resp []byte

		switch metricsType {
		case Gauge:
			resp = []byte(fmt.Sprintf("%g", *res.Value))
		case Counter:
			resp = []byte(fmt.Sprintf("%d", *res.Delta))
		default:
			http.Error(w, "metrics type is not found", http.StatusNotImplemented)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(resp)
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
