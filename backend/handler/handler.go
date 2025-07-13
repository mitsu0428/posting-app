package handler

import (
	"posting-app/domain"
)

type Handler struct {
	authUsecase         domain.AuthUsecase
	postUsecase         domain.PostUsecase
	adminUsecase        domain.AdminUsecase
	subscriptionUsecase domain.SubscriptionUsecase
}

func NewHandler(authUsecase domain.AuthUsecase, postUsecase domain.PostUsecase, adminUsecase domain.AdminUsecase, subscriptionUsecase domain.SubscriptionUsecase) *Handler {
	return &Handler{
		authUsecase:         authUsecase,
		postUsecase:         postUsecase,
		adminUsecase:        adminUsecase,
		subscriptionUsecase: subscriptionUsecase,
	}
}