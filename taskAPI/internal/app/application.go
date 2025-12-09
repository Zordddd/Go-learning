package app

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/Zordddd/learning/taskAPI/internal/http/handler"
	"github.com/Zordddd/learning/taskAPI/internal/http/middleware"
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
		middleware.AuthMiddleware,
		middleware.JsonContentTypeMiddleware,
		middleware.RequestIDMiddleware,
		rateLimiterMiddleware,
	)

	mux.HandleFunc("/tasks", chain(handler.GetTasksHandler))
	mux.HandleFunc("/tasks/create", chain(handler.GetTasksHandler))

	return mux
}
