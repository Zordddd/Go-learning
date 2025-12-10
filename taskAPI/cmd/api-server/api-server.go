package main

import (
	"log/slog"

	"github.com/Zordddd/learning/taskAPI/internal/app"
)

func main() {
	application := app.NewApplication()

	application.SetupRoutes()

	slog.Info("Server starting on :8080")
	slog.Info("Test endpoints:")
	slog.Info("  GET  /tasks (requires X-API-Key: password)")
	slog.Info("  POST /tasks/create (requires X-API-Key: password)")

	if err := application.Run(); err != nil {
		panic(err)
	}
}
