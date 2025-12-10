package middleware

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/Zordddd/learning/taskAPI/pkg/http/responseWriter"
)

func NewTimeoutMiddleware(timeout time.Duration) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()
			successHandle := make(chan struct{}, 1)
			rw := responseWriter.NewResponseWriter(w)
			go func() {
				next(rw, r.WithContext(ctx))
				successHandle <- struct{}{}
			}()

			select {
			case <-successHandle:
				return
			case <-ctx.Done():
				if !rw.Written() {
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
					if err := json.NewEncoder(rw).Encode(response); err != nil {
						http.Error(rw, err.Error(), http.StatusInternalServerError)
					}
				}
			}
		}
	}
}
