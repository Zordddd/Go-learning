package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-Api-Key")
		if apiKey != "password" {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			if err := json.NewEncoder(w).Encode(map[string]string{
				"error": "Unauthorized",
			}); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		slog.Info("Authorization user:", "status", "authorized")
		next.ServeHTTP(w, r)
	}
}
