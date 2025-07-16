package usecase

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/checkout/session"
	"github.com/stripe/stripe-go/v76/customer"
	"github.com/stripe/stripe-go/v76/subscription"
	"posting-app/domain"
	"posting-app/repository"
)

type SubscriptionUsecase struct {
	userRepo         *repository.UserRepository
	subscriptionRepo *repository.SubscriptionRepository
	stripeAPIKey     string
	priceID          string
	webhookSecret    string
	baseURL          string
}

func NewSubscriptionUsecase(
	userRepo *repository.UserRepository,
	subscriptionRepo *repository.SubscriptionRepository,
	stripeAPIKey, priceID, webhookSecret, baseURL string,
) *SubscriptionUsecase {
	stripe.Key = stripeAPIKey
	return &SubscriptionUsecase{
		userRepo:         userRepo,
		subscriptionRepo: subscriptionRepo,
		stripeAPIKey:     stripeAPIKey,
		priceID:          priceID,
		webhookSecret:    webhookSecret,
		baseURL:          baseURL,
	}
}

func (u *SubscriptionUsecase) GetSubscriptionStatus(userID int) (*domain.Subscription, error) {
	// First get the user to check the subscription status in the users table
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Try to get the subscription details from the subscriptions table
	sub, err := u.subscriptionRepo.GetByUserID(userID)
	if err != nil {
		// User has no subscription record, return status from user record
		return &domain.Subscription{
			UserID: userID,
			Status: user.SubscriptionStatus,
		}, nil
	}

	// If we have a subscription record, use that status but ensure it matches the user's status
	// The user's subscription_status should be the authoritative source
	sub.Status = user.SubscriptionStatus
	return sub, nil
}

func (u *SubscriptionUsecase) CreateCheckoutSession(userID int) (string, error) {
	user, err := u.userRepo.GetByID(userID)
	if err != nil {
		return "", errors.New("user not found")
	}

	// Create or get Stripe customer
	var customerID string
	if user.StripeCustomerID != nil {
		customerID = *user.StripeCustomerID
	} else {
		// Create new customer
		params := &stripe.CustomerParams{
			Email: stripe.String(user.Email),
			Name:  stripe.String(user.DisplayName),
		}
		c, err := customer.New(params)
		if err != nil {
			return "", fmt.Errorf("failed to create Stripe customer: %w", err)
		}
		customerID = c.ID

		// Update user with customer ID
		user.StripeCustomerID = &customerID
		err = u.userRepo.Update(user)
		if err != nil {
			slog.Error("Failed to update user with Stripe customer ID", "error", err)
		}
	}

	// Create checkout session
	params := &stripe.CheckoutSessionParams{
		Customer: stripe.String(customerID),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				Price:    stripe.String(u.priceID),
				Quantity: stripe.Int64(1),
			},
		},
		Mode:       stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		SuccessURL: stripe.String(u.baseURL + "/subscription?success=true"),
		CancelURL:  stripe.String(u.baseURL + "/subscription?canceled=true"),
	}

	s, err := session.New(params)
	if err != nil {
		return "", fmt.Errorf("failed to create checkout session: %w", err)
	}

	slog.Info("Checkout session created", "user_id", userID, "session_id", s.ID)
	return s.URL, nil
}

func (u *SubscriptionUsecase) HandleWebhook(eventType string, data interface{}) error {
	switch eventType {
	case "checkout.session.completed":
		return u.handleCheckoutSessionCompleted(data)
	case "customer.subscription.created":
		return u.handleSubscriptionCreated(data)
	case "customer.subscription.updated":
		return u.handleSubscriptionUpdated(data)
	case "customer.subscription.deleted":
		return u.handleSubscriptionDeleted(data)
	case "invoice.payment_failed":
		return u.handlePaymentFailed(data)
	default:
		slog.Info("Unhandled webhook event", "type", eventType)
		return nil
	}
}

func (u *SubscriptionUsecase) handleCheckoutSessionCompleted(data interface{}) error {
	sessionData, ok := data.(map[string]interface{})
	if !ok {
		return errors.New("invalid session data")
	}

	customerID, ok := sessionData["customer"].(string)
	if !ok {
		return errors.New("missing customer ID in session")
	}

	// Find user by Stripe customer ID
	user, err := u.findUserByStripeCustomerID(customerID)
	if err != nil {
		return fmt.Errorf("failed to find user by customer ID: %w", err)
	}

	slog.Info("Checkout session completed", "user_id", user.ID, "customer_id", customerID)
	return nil
}

func (u *SubscriptionUsecase) handleSubscriptionCreated(data interface{}) error {
	subData, ok := data.(map[string]interface{})
	if !ok {
		return errors.New("invalid subscription data")
	}

	return u.processSubscriptionEvent(subData)
}

func (u *SubscriptionUsecase) handleSubscriptionUpdated(data interface{}) error {
	subData, ok := data.(map[string]interface{})
	if !ok {
		return errors.New("invalid subscription data")
	}

	return u.processSubscriptionEvent(subData)
}

func (u *SubscriptionUsecase) handleSubscriptionDeleted(data interface{}) error {
	subData, ok := data.(map[string]interface{})
	if !ok {
		return errors.New("invalid subscription data")
	}

	return u.processSubscriptionEvent(subData)
}

func (u *SubscriptionUsecase) processSubscriptionEvent(subData map[string]interface{}) error {
	stripeSubID, ok := subData["id"].(string)
	if !ok {
		return errors.New("missing subscription ID")
	}

	customerID, ok := subData["customer"].(string)
	if !ok {
		return errors.New("missing customer ID")
	}

	status, ok := subData["status"].(string)
	if !ok {
		return errors.New("missing subscription status")
	}

	currentPeriodStart, ok := subData["current_period_start"].(float64)
	if !ok {
		return errors.New("missing current_period_start")
	}

	currentPeriodEnd, ok := subData["current_period_end"].(float64)
	if !ok {
		return errors.New("missing current_period_end")
	}

	// Find user by Stripe customer ID
	user, err := u.findUserByStripeCustomerID(customerID)
	if err != nil {
		return fmt.Errorf("failed to find user by customer ID: %w", err)
	}

	// Convert Stripe status to our domain status
	var domainStatus domain.UserSubscriptionStatus
	switch status {
	case "active":
		domainStatus = domain.UserSubscriptionStatusActive
	case "past_due":
		domainStatus = domain.UserSubscriptionStatusPastDue
	case "canceled", "incomplete_expired", "unpaid":
		domainStatus = domain.UserSubscriptionStatusCanceled
	default:
		domainStatus = domain.UserSubscriptionStatusInactive
	}

	// Update user subscription status
	user.SubscriptionStatus = domainStatus
	err = u.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update user subscription status: %w", err)
	}

	// Create or update subscription record
	existingSub, err := u.subscriptionRepo.GetByStripeSubscriptionID(stripeSubID)
	if err != nil {
		// Create new subscription
		sub := &domain.Subscription{
			UserID:               user.ID,
			StripeSubscriptionID: stripeSubID,
			Status:               domainStatus,
			CurrentPeriodStart:   time.Unix(int64(currentPeriodStart), 0),
			CurrentPeriodEnd:     time.Unix(int64(currentPeriodEnd), 0),
		}
		err = u.subscriptionRepo.Create(sub)
		if err != nil {
			return fmt.Errorf("failed to create subscription: %w", err)
		}
	} else {
		// Update existing subscription
		existingSub.Status = domainStatus
		existingSub.CurrentPeriodStart = time.Unix(int64(currentPeriodStart), 0)
		existingSub.CurrentPeriodEnd = time.Unix(int64(currentPeriodEnd), 0)
		err = u.subscriptionRepo.Update(existingSub)
		if err != nil {
			return fmt.Errorf("failed to update subscription: %w", err)
		}
	}

	slog.Info("Subscription processed", "user_id", user.ID, "status", status, "stripe_sub_id", stripeSubID)
	return nil
}

func (u *SubscriptionUsecase) handlePaymentFailed(data interface{}) error {
	invoiceData, ok := data.(map[string]interface{})
	if !ok {
		return errors.New("invalid invoice data")
	}

	customerID, ok := invoiceData["customer"].(string)
	if !ok {
		return errors.New("missing customer ID in invoice")
	}

	// Find user by Stripe customer ID
	user, err := u.findUserByStripeCustomerID(customerID)
	if err != nil {
		return fmt.Errorf("failed to find user by customer ID: %w", err)
	}

	// Update user status to past_due
	user.SubscriptionStatus = domain.UserSubscriptionStatusPastDue
	err = u.userRepo.Update(user)
	if err != nil {
		return fmt.Errorf("failed to update user subscription status: %w", err)
	}

	slog.Info("Payment failed, user status updated to past_due", "user_id", user.ID, "customer_id", customerID)
	return nil
}

func (u *SubscriptionUsecase) findUserByStripeCustomerID(customerID string) (*domain.User, error) {
	// This is a simplified implementation. In a real app, you might want to add an index on stripe_customer_id
	// or implement a more efficient query
	users, _, err := u.userRepo.GetAll(1, 1000) // Get all users (assuming we don't have too many)
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		if user.StripeCustomerID != nil && *user.StripeCustomerID == customerID {
			return user, nil
		}
	}

	return nil, errors.New("user not found with given Stripe customer ID")
}

// Batch function to sync subscription statuses
func (u *SubscriptionUsecase) SyncSubscriptionStatuses() error {
	// Get all users with Stripe customer IDs
	users, _, err := u.userRepo.GetAll(1, 1000)
	if err != nil {
		return fmt.Errorf("failed to get users: %w", err)
	}

	for _, user := range users {
		if user.StripeCustomerID == nil {
			continue
		}

		// Get Stripe subscription
		iter := subscription.List(&stripe.SubscriptionListParams{
			Customer: stripe.String(*user.StripeCustomerID),
		})

		var latestSub *stripe.Subscription
		for iter.Next() {
			sub := iter.Subscription()
			if latestSub == nil || sub.Created > latestSub.Created {
				latestSub = sub
			}
		}

		if latestSub == nil {
			// No subscription found, mark as inactive
			if user.SubscriptionStatus != domain.UserSubscriptionStatusInactive {
				user.SubscriptionStatus = domain.UserSubscriptionStatusInactive
				err = u.userRepo.Update(user)
				if err != nil {
					slog.Error("Failed to update user subscription status", "user_id", user.ID, "error", err)
				}
			}
			continue
		}

		// Convert Stripe status to our domain status
		var domainStatus domain.UserSubscriptionStatus
		switch latestSub.Status {
		case stripe.SubscriptionStatusActive:
			domainStatus = domain.UserSubscriptionStatusActive
		case stripe.SubscriptionStatusPastDue:
			domainStatus = domain.UserSubscriptionStatusPastDue
		case stripe.SubscriptionStatusCanceled, stripe.SubscriptionStatusIncompleteExpired, stripe.SubscriptionStatusUnpaid:
			domainStatus = domain.UserSubscriptionStatusCanceled
		default:
			domainStatus = domain.UserSubscriptionStatusInactive
		}

		// Update user if status changed
		if user.SubscriptionStatus != domainStatus {
			user.SubscriptionStatus = domainStatus
			err = u.userRepo.Update(user)
			if err != nil {
				slog.Error("Failed to update user subscription status", "user_id", user.ID, "error", err)
				continue
			}
		}

		// Update or create subscription record
		existingSub, err := u.subscriptionRepo.GetByUserID(user.ID)
		if err != nil {
			// Create new subscription
			sub := &domain.Subscription{
				UserID:               user.ID,
				StripeSubscriptionID: latestSub.ID,
				Status:               domainStatus,
				CurrentPeriodStart:   time.Unix(latestSub.CurrentPeriodStart, 0),
				CurrentPeriodEnd:     time.Unix(latestSub.CurrentPeriodEnd, 0),
			}
			err = u.subscriptionRepo.Create(sub)
			if err != nil {
				slog.Error("Failed to create subscription", "user_id", user.ID, "error", err)
			}
		} else {
			// Update existing subscription
			existingSub.Status = domainStatus
			existingSub.CurrentPeriodStart = time.Unix(latestSub.CurrentPeriodStart, 0)
			existingSub.CurrentPeriodEnd = time.Unix(latestSub.CurrentPeriodEnd, 0)
			err = u.subscriptionRepo.Update(existingSub)
			if err != nil {
				slog.Error("Failed to update subscription", "user_id", user.ID, "error", err)
			}
		}
	}

	slog.Info("Subscription status sync completed")
	return nil
}
