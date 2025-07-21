package di

import (
	"database/sql"

	"posting-app/handler"
	"posting-app/infrastructure"
	"posting-app/repository"
	"posting-app/usecase"
)

type Container struct {
	DB         *sql.DB
	JWTService *infrastructure.JWTService
	Handlers   *handler.Handlers
}

type Config struct {
	DB                  infrastructure.Config
	JWT                 infrastructure.JWTConfig
	StripeAPIKey        string `envconfig:"STRIPE_API_KEY" required:"true"`
	StripePriceID       string `envconfig:"STRIPE_PRICE_ID" required:"true"`
	StripeWebhookSecret string `envconfig:"STRIPE_WEBHOOK_SECRET" required:"true"`
	StripeMockMode      bool   `envconfig:"STRIPE_MOCK_MODE" default:"false"`
	BaseURL             string `envconfig:"BASE_URL" default:"http://localhost:3000"`
}

func NewContainer(config Config) (*Container, error) {
	// Database
	db, err := infrastructure.NewDatabase(config.DB)
	if err != nil {
		return nil, err
	}

	// JWT Service
	jwtService := infrastructure.NewJWTService(config.JWT)

	// Repositories
	userRepo := repository.NewUserRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)
	postRepo := repository.NewPostRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)

	// Usecases
	authUsecase := usecase.NewAuthUsecase(userRepo, passwordResetRepo, jwtService)
	postUsecase := usecase.NewPostUsecase(postRepo, userRepo)
	subscriptionUsecase := usecase.NewSubscriptionUsecase(
		userRepo,
		subscriptionRepo,
		config.StripeAPIKey,
		config.StripePriceID,
		config.StripeWebhookSecret,
		config.BaseURL,
		config.StripeMockMode,
	)

	// Handlers
	authHandler := handler.NewAuthHandler(authUsecase)
	postHandler := handler.NewPostHandler(postUsecase)
	adminHandler := handler.NewAdminHandler(authUsecase, postUsecase, userRepo)
	subscriptionHandler := handler.NewSubscriptionHandler(subscriptionUsecase, config.StripeWebhookSecret)
	userHandler := handler.NewUserHandler(userRepo)

	handlers := &handler.Handlers{
		Auth:         authHandler,
		Post:         postHandler,
		Admin:        adminHandler,
		Subscription: subscriptionHandler,
		User:         userHandler,
	}

	return &Container{
		DB:         db,
		JWTService: jwtService,
		Handlers:   handlers,
	}, nil
}
