package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)

		requestID := w.Header().Get("X-Request-ID")

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"ip", r.RemoteAddr,
			"duration", time.Since(start),
			"request_id", requestID,
		)
	}
}
