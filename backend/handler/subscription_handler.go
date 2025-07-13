package handler

import (
	"io"
	"net/http"

	"github.com/stripe/stripe-go/v76/webhook"
	"posting-app/usecase"
)

type SubscriptionHandler struct {
	subscriptionUsecase *usecase.SubscriptionUsecase
	webhookSecret       string
}

func NewSubscriptionHandler(subscriptionUsecase *usecase.SubscriptionUsecase, webhookSecret string) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionUsecase: subscriptionUsecase,
		webhookSecret:       webhookSecret,
	}
}

func (h *SubscriptionHandler) GetStatus(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	subscription, err := h.subscriptionUsecase.GetSubscriptionStatus(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	response := map[string]interface{}{
		"status": subscription.Status,
	}

	if subscription.ID != 0 {
		response["current_period_end"] = subscription.CurrentPeriodEnd
	}

	writeJSON(w, http.StatusOK, response)
}

func (h *SubscriptionHandler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	user := GetUserFromContext(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	url, err := h.subscriptionUsecase.CreateCheckoutSession(user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"url": url,
	})
}

func (h *SubscriptionHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Failed to read request body")
		return
	}

	event, err := webhook.ConstructEvent(payload, r.Header.Get("Stripe-Signature"), h.webhookSecret)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid webhook signature")
		return
	}

	err = h.subscriptionUsecase.HandleWebhook(string(event.Type), event.Data.Object)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to process webhook")
		return
	}

	w.WriteHeader(http.StatusOK)
}
