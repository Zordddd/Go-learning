package app

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zordddd/learning/taskAPI/internal/http/handler"
	"github.com/Zordddd/learning/taskAPI/internal/http/middleware"
	"github.com/Zordddd/learning/taskAPI/internal/storage"
	"github.com/Zordddd/learning/taskAPI/pkg/logger"
)

type Application struct {
	server *http.Server
	logger *slog.Logger
}

func NewApplication() *Application {
	newLogger := logger.SetupLogger()
	return &Application{
		server: &http.Server{
			Addr:         ":8080",
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
		},
		logger: newLogger,
	}
}

func (app *Application) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	rateLimiter := middleware.NewRateLimiter(time.Minute, 5)
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(rateLimiter)
	chain := middleware.Chain(
		middleware.LoggingMiddleware,
		middleware.RecoveryMiddleware,
		middleware.AuthMiddleware,
		middleware.JsonContentTypeMiddleware,
		middleware.RequestIDMiddleware,
		rateLimiterMiddleware,
	)
	mux.HandleFunc("/health", app.healthHandler)
	mux.HandleFunc("/liveness", app.livenessHandler)
	mux.HandleFunc("/readiness", app.readinessHandler)

	mux.HandleFunc("/tasks", chain(handler.GetTasksHandler))
	mux.HandleFunc("/tasks/create", chain(handler.CreateTaskHandler))

	return mux
}

func (app *Application) Run() error {
	errServer := make(chan error, 1)
	go func() {
		if err := app.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errServer <- err
		}
	}()

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errServer:
		return err
	case <-shutdown:
		shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
		defer cancel()
		if err := app.server.Shutdown(shutdownCtx); err != nil {
			slog.Warn("Failed to shutdown server gracefully", "error", err)
			return err
		}
		slog.Info("Server shutdown gracefully")
	}

	return nil
}

func (app *Application) livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (app *Application) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if bytes, err := w.Write([]byte(`{"alive": true}`)); err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{
			"error": err,
			"bytes": bytes,
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (app *Application) readinessHandler(w http.ResponseWriter, r *http.Request) {
	if storage.Database.Tasks != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"status":    "ready",
			"database":  "success init",
			"timestamp": time.Now(),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{
			"status":    "not ready",
			"database":  "bad init",
			"timestamp": time.Now(),
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
