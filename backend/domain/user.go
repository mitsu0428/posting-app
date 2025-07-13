package domain

import (
	"time"
)

type User struct {
	ID                 int                    `json:"id" db:"id"`
	Email              string                 `json:"email" db:"email"`
	PasswordHash       string                 `json:"-" db:"password_hash"`
	DisplayName        string                 `json:"display_name" db:"display_name"`
	Bio                *string                `json:"bio" db:"bio"`
	Role               string                 `json:"role" db:"role"`
	SubscriptionStatus UserSubscriptionStatus `json:"subscription_status" db:"subscription_status"`
	StripeCustomerID   *string                `json:"-" db:"stripe_customer_id"`
	IsActive           bool                   `json:"is_active" db:"is_active"`
	EmailVerified      bool                   `json:"-" db:"email_verified"`
	CreatedAt          time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at" db:"updated_at"`
}

type UserSubscriptionStatus string

const (
	UserSubscriptionStatusActive   UserSubscriptionStatus = "active"
	UserSubscriptionStatusInactive UserSubscriptionStatus = "inactive"
	UserSubscriptionStatusPastDue  UserSubscriptionStatus = "past_due"
	UserSubscriptionStatusCanceled UserSubscriptionStatus = "canceled"
)

type UserRole string

const (
	UserRoleUser  UserRole = "user"
	UserRoleAdmin UserRole = "admin"
)

type PasswordReset struct {
	ID        int        `json:"id" db:"id"`
	UserID    int        `json:"user_id" db:"user_id"`
	Token     string     `json:"token" db:"token"`
	ExpiresAt time.Time  `json:"expires_at" db:"expires_at"`
	UsedAt    *time.Time `json:"used_at" db:"used_at"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
}
