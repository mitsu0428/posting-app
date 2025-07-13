package main

import (
	"log/slog"
	"os"

	"github.com/kelseyhightower/envconfig"
	"posting-app/di"
	"posting-app/infrastructure"
	"posting-app/repository"
	"posting-app/usecase"
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

	// Initialize database
	db, err := infrastructure.NewDatabase(config.DB)
	if err != nil {
		slog.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)

	// Initialize subscription usecase
	subscriptionUsecase := usecase.NewSubscriptionUsecase(
		userRepo,
		subscriptionRepo,
		config.StripeAPIKey,
		config.StripePriceID,
		config.StripeWebhookSecret,
		config.BaseURL,
	)

	// Run subscription status sync
	slog.Info("Starting subscription status sync...")
	err = subscriptionUsecase.SyncSubscriptionStatuses()
	if err != nil {
		slog.Error("Failed to sync subscription statuses", "error", err)
		os.Exit(1)
	}

	slog.Info("Subscription status sync completed successfully")
}