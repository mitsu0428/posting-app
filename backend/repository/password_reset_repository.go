package repository

import (
	"database/sql"

	"posting-app/domain"
)

type PasswordResetRepository struct {
	db *sql.DB
}

func NewPasswordResetRepository(db *sql.DB) *PasswordResetRepository {
	return &PasswordResetRepository{db: db}
}

func (r *PasswordResetRepository) Create(reset *domain.PasswordReset) error {
	query := `
		INSERT INTO password_resets (user_id, token, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	err := r.db.QueryRow(
		query,
		reset.UserID,
		reset.Token,
		reset.ExpiresAt,
	).Scan(&reset.ID, &reset.CreatedAt)

	return err
}

func (r *PasswordResetRepository) GetByToken(token string) (*domain.PasswordReset, error) {
	reset := &domain.PasswordReset{}
	query := `
		SELECT id, user_id, token, expires_at, used_at, created_at
		FROM password_resets 
		WHERE token = $1 AND expires_at > CURRENT_TIMESTAMP AND used_at IS NULL`

	err := r.db.QueryRow(query, token).Scan(
		&reset.ID, &reset.UserID, &reset.Token, &reset.ExpiresAt, &reset.UsedAt, &reset.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return reset, nil
}

func (r *PasswordResetRepository) MarkAsUsed(id int) error {
	query := `UPDATE password_resets SET used_at = CURRENT_TIMESTAMP WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *PasswordResetRepository) DeleteExpired() error {
	query := `DELETE FROM password_resets WHERE expires_at < CURRENT_TIMESTAMP`
	_, err := r.db.Exec(query)
	return err
}
