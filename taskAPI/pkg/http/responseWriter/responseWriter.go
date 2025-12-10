package responseWriter

import "net/http"

type ResponseWriter struct {
	w       http.ResponseWriter
	written bool
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{w: w}
}

func (rw *ResponseWriter) Header() http.Header {
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

func (rw *ResponseWriter) Written() bool {
	return rw.written
}
