package middleware

import (
	"log/slog"
	"net/http"
)

func RecoveryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				slog.Error("panic", "error", err, "request", r)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				if bytes, err := w.Write([]byte(`{"error":"Internal Server Error"}`)); err != nil {
					slog.Error("Write response error", "error", err, "bytes", bytes)
				}
			}
		}()
		next(w, r)
	}
}
