package middleware

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type RateLimiter struct {
	rates    map[string]int
	mu       sync.Mutex
	limit    int
	resetAt  time.Time
	interval time.Duration
}

type RateLimiterError struct{}

func (RateLimiterError) Error() string {
	return "Rate Limiter Error"
}

func NewRateLimiter(interval time.Duration, limit int) *RateLimiter {
	return &RateLimiter{rates: map[string]int{}, limit: limit, resetAt: time.Now(), interval: interval}
}

func (rl *RateLimiter) Check(ip string) error {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	if time.Now().After(rl.resetAt) {
		rl.rates = map[string]int{}
		rl.resetAt = time.Now().Add(rl.interval)
	}

	if rl.rates[ip] >= rl.limit {
		slog.Info("IP rate limit exceeded")
		return RateLimiterError{}
	}

	rl.rates[ip]++
	return nil
}

func NewRateLimiterMiddleware(rateLimiter *RateLimiter) func(next http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ip := r.RemoteAddr
			err := rateLimiter.Check(ip)
			if err != nil {
				slog.Warn("IP rate limit exceeded", "ip", ip, "error", err)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusTooManyRequests)
				if err := json.NewEncoder(w).Encode(map[string]string{"error": "Rate limit exceeded. Try again later."}); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
				return
			}
			slog.Info("Rate limit check",
				"status", "success",
				"current_count", rateLimiter.rates[ip],
				"ip", ip,
			)
			next(w, r)
		}
	}
}
