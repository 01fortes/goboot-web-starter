package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

// Logger logs requests
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &ResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}
		next.ServeHTTP(rw, r)
		duration := time.Since(start)

		slog.Info("HTTP request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rw.StatusCode,
			"duration", duration,
			"ip", r.RemoteAddr,
			"user-agent", r.UserAgent(),
		)
	})
}
