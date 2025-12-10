package main

import (
	"log/slog"

	_ "github.com/Zordddd/learning/taskAPI/docs"
	"github.com/Zordddd/learning/taskAPI/internal/app"
)

// @title Task Management API
// @version 1.0
// @description A simple task management API with authentication and middleware
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@taskapi.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key
func main() {
	application := app.NewApplication()

	application.SetupRoutes()

	slog.Info("Server starting on :8080")
	slog.Info("Test endpoints:")
	slog.Info("  GET  /task (requires X-API-Key: password)")
	slog.Info("  POST /task (requires X-API-Key: password)")
	slog.Info("  PUT /task (requires X-API-Key: password)")
	slog.Info("  DELETE /task?id=1 (requires X-API-Key: password)")
	slog.Info("  GET /health")
	slog.Info("  GET /liveness")
	slog.Info("  GET /readiness")
	slog.Info("Swagger UI: http://localhost:8080/swagger/index.html")

	if err := application.Run(); err != nil {
		panic(err)
	}
}
