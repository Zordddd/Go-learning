package middleware

import (
	"log/slog"
	"net/http"

	"github.com/Zordddd/learning/taskAPI/pkg/http/responseWriter"
)

func JsonContentTypeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := responseWriter.NewResponseWriter(w)
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			if r.Header.Get("Content-Type") != "application/json" {
				http.Error(w, `{"error": "Content-Type must be application/json"}`, http.StatusUnsupportedMediaType)
				return
			}
		}

		slog.Info("Content-Type setup:", "status", "success")
		rw.Header().Set("Content-Type", "application/json")
		rw.WriteHeader(http.StatusOK)
		next(rw, r)
	}
}
