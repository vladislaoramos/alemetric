package server

import (
	"errors"
	"github.com/vladislaoramos/alemetric/internal/usecase"
	"net/http"
)

func errorHandler(w http.ResponseWriter, err error) {
	if errors.Is(err, usecase.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
	} else if errors.Is(err, usecase.ErrNotImplemented) {
		http.Error(w, err.Error(), http.StatusNotImplemented)
	} else {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
