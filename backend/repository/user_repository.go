package repository

import (
	"database/sql"

	"posting-app/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	query := `
		INSERT INTO users (email, password_hash, display_name, bio, role, subscription_status, is_active, email_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, created_at, updated_at`

	err := r.db.QueryRow(
		query,
		user.Email,
		user.PasswordHash,
		user.DisplayName,
		user.Bio,
		user.Role,
		user.SubscriptionStatus,
		user.IsActive,
		user.EmailVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	return err
}

func (r *UserRepository) GetByID(id int) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, password_hash, display_name, bio, role, subscription_status, 
			   stripe_customer_id, is_active, email_verified, created_at, updated_at
		FROM users WHERE id = $1 AND is_active = true`

	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName, &user.Bio,
		&user.Role, &user.SubscriptionStatus, &user.StripeCustomerID, &user.IsActive,
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, password_hash, display_name, bio, role, subscription_status, 
			   stripe_customer_id, is_active, email_verified, created_at, updated_at
		FROM users WHERE email = $1`

	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName, &user.Bio,
		&user.Role, &user.SubscriptionStatus, &user.StripeCustomerID, &user.IsActive,
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	query := `
		UPDATE users 
		SET display_name = $1, bio = $2, subscription_status = $3, stripe_customer_id = $4, 
			is_active = $5, email_verified = $6, updated_at = CURRENT_TIMESTAMP
		WHERE id = $7`

	_, err := r.db.Exec(
		query,
		user.DisplayName,
		user.Bio,
		user.SubscriptionStatus,
		user.StripeCustomerID,
		user.IsActive,
		user.EmailVerified,
		user.ID,
	)

	return err
}

func (r *UserRepository) UpdatePassword(userID int, passwordHash string) error {
	query := `UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, passwordHash, userID)
	return err
}

func (r *UserRepository) Deactivate(userID int) error {
	query := `UPDATE users SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *UserRepository) GetAll(page, limit int) ([]*domain.User, int, error) {
	offset := (page - 1) * limit

	// Get total count
	var total int
	err := r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get users
	query := `
		SELECT id, email, password_hash, display_name, bio, role, subscription_status, 
			   stripe_customer_id, is_active, email_verified, created_at, updated_at
		FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		user := &domain.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName, &user.Bio,
			&user.Role, &user.SubscriptionStatus, &user.StripeCustomerID, &user.IsActive,
			&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *UserRepository) Ban(userID int) error {
	query := `UPDATE users SET is_active = false, updated_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *UserRepository) GetByDisplayName(displayName string) (*domain.User, error) {
	user := &domain.User{}
	query := `
		SELECT id, email, password_hash, display_name, bio, role, subscription_status, 
			   stripe_customer_id, is_active, email_verified, created_at, updated_at
		FROM users WHERE display_name = $1 AND is_active = true`

	err := r.db.QueryRow(query, displayName).Scan(
		&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName, &user.Bio,
		&user.Role, &user.SubscriptionStatus, &user.StripeCustomerID, &user.IsActive,
		&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) SearchByDisplayName(query string) ([]domain.User, error) {
	sqlQuery := `
		SELECT id, email, password_hash, display_name, bio, role, subscription_status, 
			   stripe_customer_id, is_active, email_verified, created_at, updated_at
		FROM users 
		WHERE display_name ILIKE '%' || $1 || '%' AND is_active = true
		ORDER BY display_name
		LIMIT 20`

	rows, err := r.db.Query(sqlQuery, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []domain.User
	for rows.Next() {
		user := domain.User{}
		err := rows.Scan(
			&user.ID, &user.Email, &user.PasswordHash, &user.DisplayName, &user.Bio,
			&user.Role, &user.SubscriptionStatus, &user.StripeCustomerID, &user.IsActive,
			&user.EmailVerified, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
