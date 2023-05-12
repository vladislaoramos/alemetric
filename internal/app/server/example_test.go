package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi/v5"
	"github.com/vladislaoramos/alemetric/internal/entity"
	"github.com/vladislaoramos/alemetric/internal/repo"
	"github.com/vladislaoramos/alemetric/internal/usecase"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

func Example_getMetricsHandler() {
	var buf bytes.Buffer
	log := logger.New("debug", &buf)
	metricsRepo, _ := repo.NewMetricsRepo()
	tool := usecase.NewMetricsTool(metricsRepo, log)

	req, _ := http.NewRequest("GET", "/", nil)

	rec := httptest.NewRecorder()

	handler := getMetricsHandler(tool, log)
	handler(rec, req)

	fmt.Println(rec.Body.String())
}

func Example_pingHandler() {
	var buf bytes.Buffer
	log := logger.New("debug", &buf)
	metricsRepo, _ := repo.NewMetricsRepo()
	tool := usecase.NewMetricsTool(metricsRepo, log)

	req, _ := http.NewRequest("GET", "/ping", nil)

	rec := httptest.NewRecorder()

	handler := pingHandler(tool, log)
	handler(rec, req)

	fmt.Println(rec.Code)
}

func Example_updateSeveralMetricsHandler() {
	var buf bytes.Buffer
	log := logger.New("debug", &buf)
	metricsRepo, _ := repo.NewMetricsRepo()
	tool := usecase.NewMetricsTool(metricsRepo, log)

	var value entity.Gauge = 10
	reqBody := []entity.Metrics{
		{ID: "metric1", MType: "type1", Value: &value},
		{ID: "metric2", MType: "type2", Value: &value},
	}

	b, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/update", bytes.NewBuffer(b))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	updateSeveralMetricsHandler(tool, log).ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		fmt.Println("Metrics updated successfully")
	} else {
		fmt.Println("Error updating metrics")
	}

	// Output:
	// Metrics updated successfully
}

func Example_updateMetricsHandler() {
	var buf bytes.Buffer
	log := logger.New("debug", &buf)
	metricsRepo, _ := repo.NewMetricsRepo()
	tool := usecase.NewMetricsTool(metricsRepo, log)

	var value entity.Gauge = 10
	reqBody := entity.Metrics{
		ID:    "metric1",
		MType: "type1",
		Value: &value,
	}

	reqBytes, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", "/update", bytes.NewBuffer(reqBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()

	updateMetricsHandler(tool, log).ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var response entity.Metrics
		err := json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			fmt.Println("Error decoding response")
			return
		}

		fmt.Printf("Updated metrics: %+v\n", response)
	} else {
		fmt.Println("Error updating metrics")
	}

	// Output:
	// Updated metrics: {ID:metric1 MType:type1 Delta:<nil> Value:10 Hash:}
}

func Example_updateSpecificMetricsHandler() {
	var buf bytes.Buffer
	log := logger.New("debug", &buf)
	metricsRepo, _ := repo.NewMetricsRepo()
	tool := usecase.NewMetricsTool(metricsRepo, log)

	metricsType := "gauge"
	metricsName := "Alloc"
	metricsValue := "10"

	url := fmt.Sprintf("/update/%s/%s/%s", metricsType, metricsName, metricsValue)

	req, _ := http.NewRequest("POST", url, nil)

	w := httptest.NewRecorder()

	updateSpecificMetricsHandler(tool, log).ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var response entity.Metrics
		err := json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			fmt.Println("Error decoding response")
			return
		}

		fmt.Printf("Updated metrics: %+v\n", response)
	} else {
		fmt.Println("Error updating metrics")
	}

	// Output:
	// Updated metrics: {ID:metric1 MType:Gauge Delta:<nil> Value:10 Hash:}
}

func Example_getSomeMetricsHandler() {
	var buf bytes.Buffer
	log := logger.New("debug", &buf)
	metricsRepo, _ := repo.NewMetricsRepo()
	tool := usecase.NewMetricsTool(metricsRepo, log)

	metrics := entity.Metrics{
		ID:    "Alloc",
		MType: "gauge",
		Delta: nil,
		Value: nil,
	}

	body, _ := json.Marshal(metrics)

	req, _ := http.NewRequest("POST", "/metrics", bytes.NewReader(body))

	w := httptest.NewRecorder()

	getSomeMetricsHandler(tool, log).ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		var response entity.Metrics
		err := json.NewDecoder(w.Body).Decode(&response)
		if err != nil {
			fmt.Println("Error decoding response")
			return
		}

		fmt.Printf("Retrieved metrics: %+v\n", response)
	} else {
		fmt.Println("Error retrieving metrics")
	}

	// Output:
	// Retrieved metrics: {ID:metric1 MType:Gauge Delta:<nil> Value:<nil> Hash:}
}

func Example_getSpecificMetricsHandler() {
	var buf bytes.Buffer
	log := logger.New("debug", &buf)
	metricsRepo, _ := repo.NewMetricsRepo()
	tool := usecase.NewMetricsTool(metricsRepo, log)

	router := chi.NewRouter()
	router.Get("/metrics/{metricsType}/{metricsName}", getSpecificMetricsHandler(tool, log))

	req, _ := http.NewRequest("GET", "/metrics/Gauge/metric1", nil)

	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code == http.StatusOK {
		response := w.Body.String()

		fmt.Printf("Retrieved metrics: %s\n", response)
	} else {
		fmt.Println("Error retrieving metrics")
	}

	// Output:
	// Retrieved metrics: <здесь должно быть значение метрик>
}
