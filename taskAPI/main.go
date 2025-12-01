package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"
)

type Task struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Status    bool      `json:"status"`
	Timestamp time.Time `json:"timestamp"`
}

var tasks = map[int]Task{}
var nextID = 1

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

func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	result := make([]Task, 0, len(tasks))
	for _, task := range tasks {
		result = append(result, task)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	task.ID = nextID
	task.Timestamp = time.Now()
	tasks[nextID] = task
	nextID++

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(task); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func loggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		requestID := w.Header().Get("X-Request-ID")

		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"ip", r.RemoteAddr,
			"duration", time.Since(start),
			"request_id", requestID,
		)
	}
}

func authMiddleware(next http.HandlerFunc) http.HandlerFunc {
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

func jsonContentTypeMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Для POST/PUT/PATCH проверяем Content-Type
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodPatch {
			if r.Header.Get("Content-Type") != "application/json" {
				http.Error(w, `{"error": "Content-Type must be application/json"}`, http.StatusUnsupportedMediaType)
				return
			}
		}

		slog.Info("Content-Type setup:", "status", "success")
		w.Header().Set("Content-Type", "application/json")
		next(w, r)
	}
}

func newRateLimiterMiddleware(rateLimiter *RateLimiter) func(next http.HandlerFunc) http.HandlerFunc {
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

func requestIDMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		bytes := make([]byte, 16)
		if _, err := rand.Read(bytes); err != nil {
			requestID := fmt.Sprintf("req_%d", time.Now().UnixNano())
			w.Header().Set("X-Request-ID", requestID)
			slog.Warn("Failed to generate random request ID, using timestamp",
				"error", err)
		} else {
			requestID := hex.EncodeToString(bytes)
			w.Header().Set("X-Request-ID", requestID)
		}

		next(w, r)
	}
}

func main() {
	rateLimiter := NewRateLimiter(time.Minute, 5)
	rateLimiterMiddleware := newRateLimiterMiddleware(rateLimiter)
	http.HandleFunc("/tasks",
		loggingMiddleware(
			jsonContentTypeMiddleware(
				requestIDMiddleware(
					rateLimiterMiddleware(
						getTasksHandler,
					),
				),
			),
		),
	)

	http.HandleFunc("/tasks/create",
		loggingMiddleware(
			authMiddleware(
				jsonContentTypeMiddleware(
					requestIDMiddleware(
						rateLimiterMiddleware(
							createTaskHandler,
						),
					),
				),
			),
		),
	)

	slog.Info("Server starting on :8080")
	slog.Info("Test endpoints:")
	slog.Info("  GET  /tasks")
	slog.Info("  POST /tasks/create (requires X-API-Key: password)")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
