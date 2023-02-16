package server

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func gzipReadHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			errorHandler(w, err)
			return
		}
		defer gz.Close()

		r.Body = gz
		next.ServeHTTP(w, r)
	})
}

func gzipWriteHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			errorHandler(w, err)
			return
		}
		defer gz.Close()

		w.Header().Set("Content-Encoding", "gzip")
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(
			gzipWriter{
				ResponseWriter: w,
				Writer:         gz,
			},
			r,
		)
	})
}
