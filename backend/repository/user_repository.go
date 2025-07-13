package repository

import (
	"database/sql"
	"posting-app/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) domain.UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (username, email, password_hash, subscription_status, stripe_customer_id, is_admin)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at, updated_at`
	
	return r.db.QueryRow(query, user.Username, user.Email, user.PasswordHash, 
		user.SubscriptionStatus, user.StripeCustomerID, user.IsAdmin).Scan(
		&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepository) GetByID(id int) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, username, email, password_hash, subscription_status, stripe_customer_id, 
		       is_admin, is_active, created_at, updated_at
		FROM users WHERE id = $1 AND is_active = true`
	
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.SubscriptionStatus, &user.StripeCustomerID, &user.IsAdmin,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (r *userRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, username, email, password_hash, subscription_status, stripe_customer_id,
		       is_admin, is_active, created_at, updated_at
		FROM users WHERE email = $1 AND is_active = true`
	
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.SubscriptionStatus, &user.StripeCustomerID, &user.IsAdmin,
		&user.IsActive, &user.CreatedAt, &user.UpdatedAt)
	
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (r *userRepository) Update(user *domain.User) error {
	query := `
		UPDATE users 
		SET username = $2, email = $3, subscription_status = $4, 
		    stripe_customer_id = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1`
	
	_, err := r.db.Exec(query, user.ID, user.Username, user.Email,
		user.SubscriptionStatus, user.StripeCustomerID)
	return err
}

func (r *userRepository) List() ([]*domain.User, error) {
	query := `
		SELECT id, username, email, subscription_status, stripe_customer_id,
		       is_admin, is_active, created_at, updated_at
		FROM users ORDER BY created_at DESC`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.SubscriptionStatus,
			&user.StripeCustomerID, &user.IsAdmin, &user.IsActive,
			&user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *userRepository) Deactivate(id int) error {
	query := `UPDATE users SET is_active = false WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}