package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"turivo-backend/internal/domain"
)

type PasswordResetTokenRepository struct {
	db *pgxpool.Pool
}

func NewPasswordResetTokenRepository(database *pgxpool.Pool) domain.PasswordResetTokenRepository {
	return &PasswordResetTokenRepository{
		db: database,
	}
}

func (r *PasswordResetTokenRepository) Create(token *domain.PasswordResetToken) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		INSERT INTO password_reset_tokens (id, user_id, token, expires_at, used, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.Exec(ctx, query,
		token.ID,
		token.UserID,
		token.Token,
		token.ExpiresAt,
		token.Used,
		token.CreatedAt,
	)

	return err
}

func (r *PasswordResetTokenRepository) GetByToken(token string) (*domain.PasswordResetToken, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		SELECT id, user_id, token, expires_at, used, created_at
		FROM password_reset_tokens
		WHERE token = $1
	`

	var resetToken domain.PasswordResetToken
	err := r.db.QueryRow(ctx, query, token).Scan(
		&resetToken.ID,
		&resetToken.UserID,
		&resetToken.Token,
		&resetToken.ExpiresAt,
		&resetToken.Used,
		&resetToken.CreatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return &resetToken, nil
}

func (r *PasswordResetTokenRepository) MarkAsUsed(token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		UPDATE password_reset_tokens
		SET used = TRUE
		WHERE token = $1
	`

	result, err := r.db.Exec(ctx, query, token)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *PasswordResetTokenRepository) DeleteExpiredTokens() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `
		DELETE FROM password_reset_tokens
		WHERE expires_at < $1
	`

	_, err := r.db.Exec(ctx, query, time.Now())
	return err
}
