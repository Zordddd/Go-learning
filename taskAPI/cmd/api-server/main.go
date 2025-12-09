package main

import (
	"log/slog"
	"net/http"

	"github.com/Zordddd/learning/taskAPI/internal/app"
)

func main() {
	application := app.NewApplication()

	mux := application.SetupRoutes()

	slog.Info("Server starting on :8080")
	slog.Info("Test endpoints:")
	slog.Info("  GET  /tasks (requires X-API-Key: password)")
	slog.Info("  POST /tasks/create (requires X-API-Key: password)")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		panic(err)
	}
}
