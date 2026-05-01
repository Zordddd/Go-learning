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

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/Zordddd/learning/taskAPI/internal/config"
	"github.com/Zordddd/learning/taskAPI/internal/http/handler"
	"github.com/Zordddd/learning/taskAPI/internal/http/middleware"
	"github.com/Zordddd/learning/taskAPI/internal/store"
	"github.com/Zordddd/learning/taskAPI/pkg/logger"

	_ "github.com/Zordddd/learning/taskAPI/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Application struct {
	server      *http.Server
	logger      *slog.Logger
	config      middleware.CORSOptions
	repoHandler *handler.TaskRepositoryHandler
	db          *pgxpool.Pool
}

func NewApplication() *Application {
	newLogger := logger.SetupLogger()

	repoHandler, database := InitDB()

	app := &Application{
		server: &http.Server{
			Addr:         ":8080",
			WriteTimeout: time.Second * 15,
			ReadTimeout:  time.Second * 15,
			IdleTimeout:  time.Second * 60,
		},
		logger: newLogger,
		config: middleware.CORSOptions{
			Origins:     []string{"*"},
			Methods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			Headers:     []string{"Content-Type", "X-API-Key", "Authorization"},
			Credentials: true,
			MaxAge:      300,
		},
		repoHandler: repoHandler,
		db:          database,
	}

	app.server.Handler = app.SetupRoutes()

	return app
}

func InitDB() (*handler.TaskRepositoryHandler, *pgxpool.Pool) {
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, config.GetConnectionString(config.LoadConfig()))
	if err != nil {
		panic(err)
	}
	taskHandler := handler.NewTaskRepositoryHandler(store.New(pool))
	return taskHandler, pool
}

func (app *Application) PingContext(ctx context.Context) error {
	err := app.db.Ping(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (app *Application) SetupRoutes() http.Handler {
	mux := http.NewServeMux()

	rateLimiter := middleware.NewRateLimiter(time.Minute, 5)
	rateLimiterMiddleware := middleware.NewRateLimiterMiddleware(rateLimiter)
	timeoutMiddleware := middleware.NewTimeoutMiddleware(time.Second * 30)
	CORSMiddleware := middleware.NewCORSMiddleware(app.config)

	chain := middleware.Chain(
		middleware.ResponseWriterMiddleware,
		middleware.LoggingMiddleware,
		middleware.RequestIDMiddleware,
		middleware.RecoveryMiddleware,
		CORSMiddleware,
		timeoutMiddleware,
		rateLimiterMiddleware,
		//middleware.AuthMiddleware,
		middleware.JsonContentTypeMiddleware,
	)
	mux.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
	))
	mux.HandleFunc("/health", app.healthHandler)
	mux.HandleFunc("/liveness", app.livenessHandler)
	mux.HandleFunc("/readiness", app.readinessHandler)

	mux.HandleFunc("/tasks", chain(app.repoHandler.TaskHandler))

	return mux
}

func (app *Application) Run() error {
	defer func() {
		if app.db != nil {
			app.db.Close()
			app.logger.Info("Database connection closed")
		}
	}()

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

// livenessHandler godoc
// @Summary Liveness probe
// @Description Kubernetes liveness probe endpoint
// @Tags health
// @Success 200
// @Router /liveness [get]
func (app *Application) livenessHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

// healthHandler godoc
// @Summary Health check endpoint
// @Description Check if the service is alive
// @Tags health
// @Produce json
// @Success 200 {object} map[string]bool
// @Failure 500 {object} map[string]interface{}
// @Router /health [get]
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

// readinessHandler godoc
// @Summary Readiness probe
// @Description Kubernetes readiness probe endpoint
// @Tags health
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /readiness [get]
func (app *Application) readinessHandler(w http.ResponseWriter, r *http.Request) {
	if app.db != nil {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()
		err := app.PingContext(ctx)
		if err != nil {
			slog.Warn("Failed to ping database", "error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			response := map[string]interface{}{
				"status":    "not ready",
				"database":  "bad init",
				"timestamp": time.Now(),
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"status":    "ready",
			"database":  "success init",
			"timestamp": time.Now(),
		}
		if err = json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusServiceUnavailable)
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
