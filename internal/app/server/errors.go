package server

import (
	"errors"
	"net/http"

	"github.com/vladislaoramos/alemetric/internal/usecase"
)

func errorHandler(w http.ResponseWriter, err error) {
	if errors.Is(err, usecase.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else if errors.Is(err, usecase.ErrNotImplemented) {
		http.Error(w, err.Error(), http.StatusNotImplemented)
	} else if errors.Is(err, usecase.ErrDataSignNotEqual) {
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
