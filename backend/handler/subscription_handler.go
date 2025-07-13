package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"posting-app/domain"
)

func (h *Handler) CreateCheckoutSession(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(UserContextKey).(*domain.User)

	sessionID, err := h.subscriptionUsecase.CreateCheckoutSession(user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := map[string]string{
		"session_id": sessionID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func (h *Handler) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	signature := r.Header.Get("Stripe-Signature")
	if signature == "" {
		http.Error(w, "Missing Stripe-Signature header", http.StatusBadRequest)
		return
	}

	if err := h.subscriptionUsecase.HandleWebhook(payload, signature); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
}