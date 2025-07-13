package di

import (
	"database/sql"
	"posting-app/domain"
	"posting-app/handler"
	"posting-app/infrastructure"
	"posting-app/repository"
	"posting-app/usecase"

	"go.uber.org/dig"
)

func BuildContainer() *dig.Container {
	container := dig.New()

	container.Provide(infrastructure.NewDatabase)
	container.Provide(infrastructure.NewJWTManager)

	container.Provide(repository.NewUserRepository)
	container.Provide(repository.NewPostRepository)
	container.Provide(repository.NewReplyRepository)
	container.Provide(repository.NewSubscriptionRepository)
	container.Provide(repository.NewPasswordResetRepository)

	container.Provide(usecase.NewAuthUsecase)
	container.Provide(usecase.NewPostUsecase)
	container.Provide(usecase.NewAdminUsecase)
	container.Provide(usecase.NewSubscriptionUsecase)

	container.Provide(handler.NewHandler)

	return container
}

type DIParams struct {
	dig.In
	Database     *sql.DB
	Handler      *handler.Handler
	UserRepo     domain.UserRepository
	PostRepo     domain.PostRepository
	ReplyRepo    domain.ReplyRepository
	AuthUsecase  domain.AuthUsecase
	PostUsecase  domain.PostUsecase
	AdminUsecase domain.AdminUsecase
}