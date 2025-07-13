package main

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
	"posting-app/di"
	"posting-app/handler"
)

func main() {
	// Setup logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Load configuration
	var config di.Config
	err := envconfig.Process("", &config)
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize container
	container, err := di.NewContainer(config)
	if err != nil {
		slog.Error("Failed to initialize container", "error", err)
		os.Exit(1)
	}
	defer container.DB.Close()

	// Setup router
	router := handler.NewRouter(container.Handlers, container.JWTService)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	slog.Info("Starting server", "port", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		slog.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}
