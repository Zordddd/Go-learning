package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/Zordddd/learning/taskAPI/pkg/http/responseWriter"
)

func RequestIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes := make([]byte, 16)
		rw := responseWriter.NewResponseWriter(w)
		if _, err := rand.Read(bytes); err != nil {
			requestID := fmt.Sprintf("req_%d", time.Now().UnixNano())
			rw.Header().Set("X-Request-ID", requestID)
			slog.Warn("Failed to generate random request ID, using timestamp",
				"error", err)
		} else {
			requestID := hex.EncodeToString(bytes)
			rw.Header().Set("X-Request-ID", requestID)
		}

		next(rw, r)
	}
}
