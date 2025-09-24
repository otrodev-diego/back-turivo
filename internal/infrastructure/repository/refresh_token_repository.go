package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"turivo-backend/internal/domain"
	"turivo-backend/internal/infrastructure/db/sqlc"
)

type RefreshTokenRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewRefreshTokenRepository(db *pgxpool.Pool) *RefreshTokenRepository {
	return &RefreshTokenRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *RefreshTokenRepository) Create(token *domain.RefreshToken) error {
	ctx := context.Background()

	dbToken, err := r.queries.CreateRefreshToken(ctx, sqlc.CreateRefreshTokenParams{
		ID:        pgtype.UUID{Bytes: token.ID, Valid: true},
		UserID:    pgtype.UUID{Bytes: token.UserID, Valid: true},
		Token:     token.Token,
		ExpiresAt: pgtype.Timestamptz{Time: token.ExpiresAt, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("failed to create refresh token: %w", err)
	}

	// Update token with created timestamp
	token.CreatedAt = dbToken.CreatedAt.Time

	return nil
}

func (r *RefreshTokenRepository) GetByToken(token string) (*domain.RefreshToken, error) {
	ctx := context.Background()

	dbToken, err := r.queries.GetRefreshTokenByToken(ctx, token)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrRefreshTokenNotFound
		}
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	var tokenID, userID uuid.UUID
	copy(tokenID[:], dbToken.ID.Bytes[:])
	copy(userID[:], dbToken.UserID.Bytes[:])

	return &domain.RefreshToken{
		ID:        tokenID,
		UserID:    userID,
		Token:     dbToken.Token,
		ExpiresAt: dbToken.ExpiresAt.Time,
		CreatedAt: dbToken.CreatedAt.Time,
	}, nil
}

func (r *RefreshTokenRepository) Delete(token string) error {
	ctx := context.Background()

	err := r.queries.DeleteRefreshToken(ctx, token)
	if err != nil {
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	return nil
}

func (r *RefreshTokenRepository) DeleteByUserID(userID uuid.UUID) error {
	ctx := context.Background()

	err := r.queries.DeleteRefreshTokensByUserID(ctx, pgtype.UUID{Bytes: userID, Valid: true})
	if err != nil {
		return fmt.Errorf("failed to delete refresh tokens by user ID: %w", err)
	}

	return nil
}
