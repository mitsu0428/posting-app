package repository

import (
	"database/sql"
	"posting-app/domain"
)

type passwordResetRepository struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) domain.PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

func (r *passwordResetRepository) Create(reset *domain.PasswordReset) error {
	query := `
		INSERT INTO password_resets (user_id, token, expires_at, created_at)
		VALUES ($1, $2, $3, NOW())
		RETURNING id, created_at
	`
	return r.db.QueryRow(query, reset.UserID, reset.Token, reset.ExpiresAt).
		Scan(&reset.ID, &reset.CreatedAt)
}

func (r *passwordResetRepository) GetByToken(token string) (*domain.PasswordReset, error) {
	query := `
		SELECT id, user_id, token, expires_at, created_at
		FROM password_resets
		WHERE token = $1 AND expires_at > NOW()
	`
	reset := &domain.PasswordReset{}
	err := r.db.QueryRow(query, token).Scan(
		&reset.ID,
		&reset.UserID,
		&reset.Token,
		&reset.ExpiresAt,
		&reset.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return reset, nil
}

func (r *passwordResetRepository) DeleteByUserID(userID int) error {
	query := `DELETE FROM password_resets WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *passwordResetRepository) DeleteExpired() error {
	query := `DELETE FROM password_resets WHERE expires_at <= NOW()`
	_, err := r.db.Exec(query)
	return err
}