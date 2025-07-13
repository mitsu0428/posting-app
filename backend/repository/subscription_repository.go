package repository

import (
	"database/sql"

	"posting-app/domain"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(subscription *domain.Subscription) error {
	query := `
		INSERT INTO subscriptions (user_id, stripe_subscription_id, status, current_period_start, current_period_end)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		subscription.UserID,
		subscription.StripeSubscriptionID,
		subscription.Status,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
	).Scan(&subscription.ID, &subscription.CreatedAt, &subscription.UpdatedAt)

	return err
}

func (r *SubscriptionRepository) GetByUserID(userID int) (*domain.Subscription, error) {
	subscription := &domain.Subscription{}
	query := `
		SELECT id, user_id, stripe_subscription_id, status, current_period_start, current_period_end, created_at, updated_at
		FROM subscriptions 
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1`

	err := r.db.QueryRow(query, userID).Scan(
		&subscription.ID, &subscription.UserID, &subscription.StripeSubscriptionID,
		&subscription.Status, &subscription.CurrentPeriodStart, &subscription.CurrentPeriodEnd,
		&subscription.CreatedAt, &subscription.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (r *SubscriptionRepository) GetByStripeSubscriptionID(stripeSubID string) (*domain.Subscription, error) {
	subscription := &domain.Subscription{}
	query := `
		SELECT id, user_id, stripe_subscription_id, status, current_period_start, current_period_end, created_at, updated_at
		FROM subscriptions 
		WHERE stripe_subscription_id = $1`

	err := r.db.QueryRow(query, stripeSubID).Scan(
		&subscription.ID, &subscription.UserID, &subscription.StripeSubscriptionID,
		&subscription.Status, &subscription.CurrentPeriodStart, &subscription.CurrentPeriodEnd,
		&subscription.CreatedAt, &subscription.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return subscription, nil
}

func (r *SubscriptionRepository) Update(subscription *domain.Subscription) error {
	query := `
		UPDATE subscriptions 
		SET status = $1, current_period_start = $2, current_period_end = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $4`

	_, err := r.db.Exec(
		query,
		subscription.Status,
		subscription.CurrentPeriodStart,
		subscription.CurrentPeriodEnd,
		subscription.ID,
	)

	return err
}
