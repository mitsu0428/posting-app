package domain

import "time"

type User struct {
	ID                 int       `json:"id" db:"id"`
	Username           string    `json:"username" db:"username"`
	Email              string    `json:"email" db:"email"`
	PasswordHash       string    `json:"-" db:"password_hash"`
	SubscriptionStatus string    `json:"subscription_status" db:"subscription_status"`
	StripeCustomerID   *string   `json:"stripe_customer_id,omitempty" db:"stripe_customer_id"`
	IsAdmin            bool      `json:"is_admin" db:"is_admin"`
	IsActive           bool      `json:"is_active" db:"is_active"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id int) (*User, error)
	GetByEmail(email string) (*User, error)
	Update(user *User) error
	List() ([]*User, error)
	Deactivate(id int) error
}

type AuthUsecase interface {
	Register(username, email, password string) (*User, error)
	Login(email, password string) (string, *User, error)
	AdminLogin(email, password string) (string, *User, error)
	ValidateToken(token string) (*User, error)
	Logout(token string) error
	ForgotPassword(email string) error
	ResetPassword(token, newPassword string) error
}

type PasswordReset struct {
	ID        int       `json:"id" db:"id"`
	UserID    int       `json:"user_id" db:"user_id"`
	Token     string    `json:"token" db:"token"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type PasswordResetRepository interface {
	Create(reset *PasswordReset) error
	GetByToken(token string) (*PasswordReset, error)
	DeleteByUserID(userID int) error
	DeleteExpired() error
}