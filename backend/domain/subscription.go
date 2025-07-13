package domain

import "time"

type Subscription struct {
	ID                   int       `json:"id" db:"id"`
	UserID               int       `json:"user_id" db:"user_id"`
	StripeSubscriptionID string    `json:"stripe_subscription_id" db:"stripe_subscription_id"`
	Status               string    `json:"status" db:"status"`
	CurrentPeriodStart   time.Time `json:"current_period_start" db:"current_period_start"`
	CurrentPeriodEnd     time.Time `json:"current_period_end" db:"current_period_end"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

type SubscriptionRepository interface {
	Create(subscription *Subscription) error
	GetByUserID(userID int) (*Subscription, error)
	GetByStripeID(stripeSubscriptionID string) (*Subscription, error)
	Update(subscription *Subscription) error
	List() ([]*Subscription, error)
}

type SubscriptionUsecase interface {
	CreateCheckoutSession(userID int) (string, error)
	HandleWebhook(payload []byte, signature string) error
}