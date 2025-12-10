package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/Zordddd/learning/taskAPI/pkg/http/responseWriter"
)

type CORSOptions struct {
	Origins     []string
	Methods     []string
	Headers     []string
	Credentials bool
	MaxAge      int
}

func NewCORSMiddleware(config CORSOptions) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			rw := responseWriter.NewResponseWriter(w)
			origin := r.Header.Get("Origin")
			if len(config.Origins) > 0 {
				allowed := false
				for _, o := range config.Origins {
					if o == origin || o == "*" {
						allowed = true
						break
					}
				}
				if !allowed {
					rw.WriteHeader(http.StatusForbidden)
					return
				}
				rw.Header().Set("Access-Control-Allow-Origin", origin)
			}
			if config.Credentials {
				rw.Header().Set("Access-Control-Allow-Credentials", "true")
			}

			// Preflight
			if r.Method == http.MethodOptions {
				if len(config.Headers) > 0 {
					rw.Header().Set("Access-Control-Allow-Headers", strings.Join(config.Headers, ","))
				}
				if len(config.Methods) > 0 {
					rw.Header().Set("Access-Control-Allow-Methods", strings.Join(config.Methods, ","))
				}
				if config.MaxAge > 0 {
					rw.Header().Set("Access-Control-Max-Age", strconv.Itoa(config.MaxAge))
				}
			}

			next(rw, r)
		}
	}
}
