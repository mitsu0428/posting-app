package usecase

import (
	"encoding/json"
	"fmt"
	"os"
	"posting-app/domain"
	"time"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/webhook"
)

type subscriptionUsecase struct {
	subscriptionRepo domain.SubscriptionRepository
	userRepo         domain.UserRepository
}

func NewSubscriptionUsecase(subscriptionRepo domain.SubscriptionRepository, userRepo domain.UserRepository) domain.SubscriptionUsecase {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")
	return &subscriptionUsecase{
		subscriptionRepo: subscriptionRepo,
		userRepo:         userRepo,
	}
}

func (u *subscriptionUsecase) CreateCheckoutSession(userID int) (string, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	priceID := os.Getenv("STRIPE_PRICE_ID")
	if priceID == "" {
		return "", fmt.Errorf("stripe price ID not configured")
	}

	params := &stripe.CheckoutSessionParams{
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(os.Getenv("FRONTEND_URL") + "/subscription/success"),
		CancelURL:  stripe.String(os.Getenv("FRONTEND_URL") + "/subscription/cancel"),
		Metadata: map[string]string{
			"user_id": fmt.Sprintf("%d", userID),
		},
	}

	if user.StripeCustomerID != nil {
		params.Customer = user.StripeCustomerID
	} else {
		params.CustomerEmail = stripe.String(user.Email)
	}

	s, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	return s.ID, nil
}

func (u *subscriptionUsecase) HandleWebhook(payload []byte, signature string) error {
	endpointSecret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if endpointSecret == "" {
		return fmt.Errorf("stripe webhook secret not configured")
	}

	event, err := webhook.ConstructEvent(payload, signature, endpointSecret)
	if err != nil {
		return fmt.Errorf("failed to verify webhook signature: %w", err)
	}

	switch event.Type {
	case "checkout.session.completed":
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return fmt.Errorf("failed to parse session: %w", err)
		}
		
		return u.handleCheckoutSessionCompleted(&session)

	case "customer.subscription.created":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			return fmt.Errorf("failed to parse subscription: %w", err)
		}
		
		return u.handleSubscriptionCreated(&subscription)

	case "customer.subscription.updated":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			return fmt.Errorf("failed to parse subscription: %w", err)
		}
		
		return u.handleSubscriptionUpdated(&subscription)

	case "customer.subscription.deleted":
		var subscription stripe.Subscription
		if err := json.Unmarshal(event.Data.Raw, &subscription); err != nil {
			return fmt.Errorf("failed to parse subscription: %w", err)
		}
		
		return u.handleSubscriptionDeleted(&subscription)
	}

	return nil
}

func (u *subscriptionUsecase) handleCheckoutSessionCompleted(session *stripe.CheckoutSession) error {
	userID := session.Metadata["user_id"]
	if userID == "" {
		return fmt.Errorf("user_id not found in session metadata")
	}

	user, err := u.userRepo.GetByID(parseInt(userID))
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if user.StripeCustomerID == nil {
		customerID := session.Customer.ID
		user.StripeCustomerID = &customerID
		if err := u.userRepo.Update(user); err != nil {
			return fmt.Errorf("failed to update user with customer ID: %w", err)
		}
	}

	return nil
}

func (u *subscriptionUsecase) handleSubscriptionCreated(stripeSubscription *stripe.Subscription) error {
	userID, err := u.getUserIDFromCustomer(stripeSubscription.Customer.ID)
	if err != nil {
		return err
	}

	subscription := &domain.Subscription{
		UserID:               userID,
		StripeSubscriptionID: stripeSubscription.ID,
		Status:               string(stripeSubscription.Status),
		CurrentPeriodStart:   time.Unix(stripeSubscription.CurrentPeriodStart, 0),
		CurrentPeriodEnd:     time.Unix(stripeSubscription.CurrentPeriodEnd, 0),
	}

	if err := u.subscriptionRepo.Create(subscription); err != nil {
		return fmt.Errorf("failed to create subscription: %w", err)
	}

	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.SubscriptionStatus = "active"
	if err := u.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user subscription status: %w", err)
	}

	return nil
}

func (u *subscriptionUsecase) handleSubscriptionUpdated(stripeSubscription *stripe.Subscription) error {
	subscription, err := u.subscriptionRepo.GetByStripeID(stripeSubscription.ID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	subscription.Status = string(stripeSubscription.Status)
	subscription.CurrentPeriodStart = time.Unix(stripeSubscription.CurrentPeriodStart, 0)
	subscription.CurrentPeriodEnd = time.Unix(stripeSubscription.CurrentPeriodEnd, 0)

	if err := u.subscriptionRepo.Update(subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	user, err := u.userRepo.GetByID(subscription.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	if stripeSubscription.Status == stripe.SubscriptionStatusActive {
		user.SubscriptionStatus = "active"
	} else {
		user.SubscriptionStatus = "inactive"
	}

	if err := u.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user subscription status: %w", err)
	}

	return nil
}

func (u *subscriptionUsecase) handleSubscriptionDeleted(stripeSubscription *stripe.Subscription) error {
	subscription, err := u.subscriptionRepo.GetByStripeID(stripeSubscription.ID)
	if err != nil {
		return fmt.Errorf("subscription not found: %w", err)
	}

	subscription.Status = "canceled"
	if err := u.subscriptionRepo.Update(subscription); err != nil {
		return fmt.Errorf("failed to update subscription: %w", err)
	}

	user, err := u.userRepo.GetByID(subscription.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	user.SubscriptionStatus = "canceled"
	if err := u.userRepo.Update(user); err != nil {
		return fmt.Errorf("failed to update user subscription status: %w", err)
	}

	return nil
}

func (u *subscriptionUsecase) getUserIDFromCustomer(customerID string) (int, error) {
	// This would typically involve querying the users table by stripe_customer_id
	// For now, return an error as this requires additional implementation
	return 0, fmt.Errorf("not implemented: getUserIDFromCustomer")
}

func parseInt(s string) int {
	// Simple implementation for demo - in production, use strconv.Atoi with error handling
	return 0
}