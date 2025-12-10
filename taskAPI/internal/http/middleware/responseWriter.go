package middleware

import (
	"net/http"

	"github.com/Zordddd/learning/taskAPI/pkg/http/responseWriter"
)

func ResponseWriterMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rw := responseWriter.NewResponseWriter(w)
		next(rw, r)
	}
}
