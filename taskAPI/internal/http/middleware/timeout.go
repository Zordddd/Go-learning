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
			rw := NewResponseWriter(w)
			go func() {
				next(rw, r.WithContext(ctx))
				successHandle <- struct{}{}
			}()

			select {
			case <-successHandle:
				return
			case <-ctx.Done():
				if !rw.written {
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

type ResponseWriter struct {
	w       http.ResponseWriter
	written bool
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w: w}
}

func (rw *ResponseWriter) Header() http.Header {
	rw.written = true
	return rw.w.Header()
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.written = true
	return rw.w.Write(b)
}

func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.written = true
	rw.w.WriteHeader(statusCode)
}
