package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Zordddd/learning/taskAPI/pkg/http/responseWriter"
)

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := responseWriter.NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		requestID := rw.Header().Get("X-Request-ID")

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"ip", r.RemoteAddr,
			"duration", time.Since(start),
			"request_id", requestID,
		)
	}
}
