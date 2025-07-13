package repository

import (
	"database/sql"
	"posting-app/domain"
)

type subscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) domain.SubscriptionRepository {
	return &subscriptionRepository{db: db}
}

func (r *subscriptionRepository) Create(subscription *domain.Subscription) error {
	query := `
		INSERT INTO subscriptions (user_id, stripe_subscription_id, status, current_period_start, current_period_end)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at`
	
	return r.db.QueryRow(query, subscription.UserID, subscription.StripeSubscriptionID,
		subscription.Status, subscription.CurrentPeriodStart, subscription.CurrentPeriodEnd).Scan(
		&subscription.ID, &subscription.CreatedAt, &subscription.UpdatedAt)
}

func (r *subscriptionRepository) GetByUserID(userID int) (*domain.Subscription, error) {
	subscription := &domain.Subscription{}
	query := `
		SELECT id, user_id, stripe_subscription_id, status, current_period_start, current_period_end, created_at, updated_at
		FROM subscriptions WHERE user_id = $1`
	
	err := r.db.QueryRow(query, userID).Scan(
		&subscription.ID, &subscription.UserID, &subscription.StripeSubscriptionID,
		&subscription.Status, &subscription.CurrentPeriodStart, &subscription.CurrentPeriodEnd,
		&subscription.CreatedAt, &subscription.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return subscription, nil
}

func (r *subscriptionRepository) GetByStripeID(stripeSubscriptionID string) (*domain.Subscription, error) {
	subscription := &domain.Subscription{}
	query := `
		SELECT id, user_id, stripe_subscription_id, status, current_period_start, current_period_end, created_at, updated_at
		FROM subscriptions WHERE stripe_subscription_id = $1`
	
	err := r.db.QueryRow(query, stripeSubscriptionID).Scan(
		&subscription.ID, &subscription.UserID, &subscription.StripeSubscriptionID,
		&subscription.Status, &subscription.CurrentPeriodStart, &subscription.CurrentPeriodEnd,
		&subscription.CreatedAt, &subscription.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return subscription, nil
}

func (r *subscriptionRepository) Update(subscription *domain.Subscription) error {
	query := `
		UPDATE subscriptions 
		SET status = $2, current_period_start = $3, current_period_end = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`
	
	_, err := r.db.Exec(query, subscription.ID, subscription.Status,
		subscription.CurrentPeriodStart, subscription.CurrentPeriodEnd)
	return err
}

func (r *subscriptionRepository) List() ([]*domain.Subscription, error) {
	query := `
		SELECT id, user_id, stripe_subscription_id, status, current_period_start, current_period_end, created_at, updated_at
		FROM subscriptions ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*domain.Subscription
	for rows.Next() {
		subscription := &domain.Subscription{}
		err := rows.Scan(
			&subscription.ID, &subscription.UserID, &subscription.StripeSubscriptionID,
			&subscription.Status, &subscription.CurrentPeriodStart, &subscription.CurrentPeriodEnd,
			&subscription.CreatedAt, &subscription.UpdatedAt)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, subscription)
	}

	return subscriptions, nil
}