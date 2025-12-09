package main

import (
	"log/slog"
	"net/http"

	"github.com/Zordddd/learning/taskAPI/internal/app"
)

func main() {
	application := app.NewApplication()

	mux := application.SetupRoutes()

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}

	slog.Info("Server starting on :8080")
	slog.Info("Test endpoints:")
	slog.Info("  GET  /tasks")
	slog.Info("  POST /tasks/create (requires X-API-Key: password)")
}
