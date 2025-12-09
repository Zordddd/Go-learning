package middleware

import (
	"log/slog"
	"net/http"
)

func JsonContentTypeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			if r.Header.Get("Content-Type") != "application/json" {
				http.Error(w, `{"error": "Content-Type must be application/json"}`, http.StatusUnsupportedMediaType)
				return
			}
		}

		slog.Info("Content-Type setup:", "status", "success")
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}
