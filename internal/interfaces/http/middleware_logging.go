package http

import (
	"log"
	"net/http"
	"time"
)

// responseRecorder wraps http.ResponseWriter to capture the status code
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

// RequestLogger logs incoming HTTP requests with method, path, status, and duration
type RequestLogger struct{}

// NewRequestLogger creates a new request logging middleware
func NewRequestLogger() *RequestLogger {
	return &RequestLogger{}
}

// Wrap wraps a handler with structured request logging
func (rl *RequestLogger) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &responseRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(rec, r)

		duration := time.Since(start)

		log.Printf("[%s] %s %s | %d | %s | %s",
			r.Method,
			r.URL.Path,
			r.Proto,
			rec.statusCode,
			duration.Round(time.Microsecond),
			r.RemoteAddr,
		)
	})
}
