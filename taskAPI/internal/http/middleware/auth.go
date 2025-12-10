package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Zordddd/learning/taskAPI/pkg/http/responseWriter"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := responseWriter.NewResponseWriter(w)
		apiKey := r.Header.Get("X-Api-Key")
		if apiKey != "password" {
			rw.Header().Set("Content-type", "application/json")
			rw.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error": "Unauthorized",
			}); err != nil {
				http.Error(rw, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		slog.Info("Authorization user:", "status", "authorized")
		next.ServeHTTP(rw, r)
	}
}
