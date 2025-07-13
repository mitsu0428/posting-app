package main

import (
	"log"
	"os"
	"posting-app/di"
	"posting-app/domain"
	"time"

	"github.com/stripe/stripe-go/v76"
	"github.com/stripe/stripe-go/v76/subscription"
)

func main() {
	stripe.Key = os.Getenv("STRIPE_SECRET_KEY")

	container := di.BuildContainer()

	err := container.Invoke(func(params di.DIParams) error {
		return runSubscriptionCheck(params.UserRepo)
	})

	if err != nil {
		log.Fatal("Failed to run subscription batch:", err)
	}
}

func runSubscriptionCheck(userRepo domain.UserRepository) error {
	log.Println("Starting subscription status check...")

	users, err := userRepo.List()
	if err != nil {
		return err
	}

	for _, user := range users {
		if user.StripeCustomerID == nil || user.SubscriptionStatus == "inactive" {
			continue
		}

		if err := checkUserSubscription(user, userRepo); err != nil {
			log.Printf("Error checking subscription for user %d: %v", user.ID, err)
		}
	}

	log.Println("Subscription status check completed")
	return nil
}

func checkUserSubscription(user *domain.User, userRepo domain.UserRepository) error {
	params := &stripe.SubscriptionListParams{
		Customer: user.StripeCustomerID,
		Status:   stripe.String("all"),
	}
	
	i := subscription.List(params)
	
	activeSubscription := false
	for i.Next() {
		subscription := i.Subscription()
		if subscription.Status == stripe.SubscriptionStatusActive {
			if time.Unix(subscription.CurrentPeriodEnd, 0).After(time.Now()) {
				activeSubscription = true
				break
			}
		}
	}

	if i.Err() != nil {
		return i.Err()
	}

	newStatus := "inactive"
	if activeSubscription {
		newStatus = "active"
	}

	if user.SubscriptionStatus != newStatus {
		user.SubscriptionStatus = newStatus
		if err := userRepo.Update(user); err != nil {
			return err
		}
		log.Printf("Updated user %d subscription status to %s", user.ID, newStatus)
	}

	return nil
}