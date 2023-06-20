package server

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime/debug"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/vladislaoramos/alemetric/internal/usecase"
	logger "github.com/vladislaoramos/alemetric/pkg/log"
)

const (
	Gauge   = "gauge"
	Counter = "counter"
)

func NewRouter(
	handler *chi.Mux,
	tool *usecase.ToolUseCase,
	l logger.LogInterface,
	privateKeyPath, cidr string,
) {
	handler.Use(middleware.RequestID)
	handler.Use(middleware.RealIP)
	handler.Use(middleware.Logger)
	handler.Use(middleware.Compress(1, "gzip"))
	handler.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					l.Error(fmt.Sprintf("panic; stacktrace: %s", string(debug.Stack())))
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()

			next.ServeHTTP(w, r)
		})
	})

	handler.Use(gzipWriteHandler)
	handler.Use(gzipReadHandler)
	handler.Use(rsaHandler(privateKeyPath))
	handler.Use(trustedSubnet(cidr))

	handler.Get("/ping", pingHandler(tool, l))

	handler.Get("/", getMetricsHandler(tool, l))

	// update
	handler.Post("/updates/", updateSeveralMetricsHandler(tool, l))
	handler.Route("/update", func(r chi.Router) {
		r.Post("/", updateMetricsHandler(tool, l))
		r.Post("/{metricsType}/{metricsName}/{metricsValue}", updateSpecificMetricsHandler(tool, l))
	})

	// value
	handler.Route("/value", func(r chi.Router) {
		r.Post("/", getSomeMetricsHandler(tool, l))
		r.Get("/{metricsType}/{metricsName}", getSpecificMetricsHandler(tool, l))
	})

	handler.Route("/debug/pprof/", func(r chi.Router) {
		r.Get("/", pprof.Index)
		r.Get("/profile", pprof.Profile)
		r.Get("/cmdline", pprof.Cmdline)
		r.Get("/symbol", pprof.Symbol)
		r.Get("/trace", pprof.Trace)
		r.Get("/{cmd}", pprof.Index)
	})
}
