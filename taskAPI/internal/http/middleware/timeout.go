package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"
)

func NewTimeoutMiddleware(timeout time.Duration) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			successHandle := make(chan struct{}, 1)
			go func() {
				defer func() {
					successHandle <- struct{}{}
				}()
				next(w, r.WithContext(ctx))
			}()

			select {
			case <-successHandle:
				return
			case <-ctx.Done():
				if rw, ok := w.(interface{ Written() bool }); ok && !rw.Written() {
					slog.Warn("timeout error",
						"timeout", timeout,
						"Addr", r.RemoteAddr,
						"path", r.URL.Path,
					)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusRequestTimeout)
					response := map[string]interface{}{
						"timeout error": timeout,
					}
					if err := json.NewEncoder(w).Encode(response); err != nil {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				}
			}
		}
	}
}
