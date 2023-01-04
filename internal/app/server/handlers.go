package server

import (
	"github.com/vladislaoramos/alemetric/internal/entity"
	"net/http"
	"strings"
)

func UpdateMetricsHandler(repo MetricsRepo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		params := strings.Split(
			strings.TrimPrefix(r.URL.Path, "/update/"),
			"/",
		)
		if len(params) != 3 {
			http.Error(w, "Required pattern is type/name/value", http.StatusNotFound)
			return
		}

		switch params[0] {
		case "counter":
			value, err := entity.ParseCounterMetrics(params[2])
			if err != nil {
				http.Error(w, "Unsuitable value type", http.StatusBadRequest)
				return
			}
			err = repo.StoreCounterMetrics(params[1], value)
			if err != nil {
				http.Error(w, "Some problem with the storage", http.StatusInternalServerError)
			}
		case "gauge":
			value, err := entity.ParseGaugeMetrics(params[2])
			if err != nil {
				http.Error(w, "Unsuitable value type", http.StatusBadRequest)
				return
			}
			err = repo.StoreGaugeMetrics(params[1], value)
			if err != nil {
				http.Error(w, "Some problem with the storage", http.StatusInternalServerError)
			}
		default:
			http.Error(w, "Metrics type is not found", http.StatusNotImplemented)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
