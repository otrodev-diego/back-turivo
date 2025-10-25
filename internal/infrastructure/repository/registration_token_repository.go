package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type RegistrationTokenRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewRegistrationTokenRepository(db *sql.DB, logger *zap.Logger) *RegistrationTokenRepository {
	return &RegistrationTokenRepository{
		db:     db,
		logger: logger,
	}
}

func (r *RegistrationTokenRepository) Create(token *domain.RegistrationToken) error {
	r.logger.Info("🗄️ === RegistrationTokenRepository.Create Started ===",
		zap.String("token_id", token.ID.String()),
		zap.String("email", token.Email),
		zap.String("role", string(token.Role)),
	)

	ctx := context.Background()

	orgID := pgtype.UUID{Valid: false}
	if token.OrgID != nil {
		orgID = pgtype.UUID{
			Bytes: *token.OrgID,
			Valid: true,
		}
		r.logger.Info("🏢 org_id provided", zap.String("org_id", token.OrgID.String()))
	} else {
		r.logger.Info("ℹ️ No org_id provided")
	}

	var companyProfile *string
	if token.CompanyProfile != nil {
		profile := string(*token.CompanyProfile)
		companyProfile = &profile
		r.logger.Info("👔 company_profile provided", zap.String("company_profile", profile))
	}

	// Since we don't have SQLC queries for registration tokens yet,
	// we'll use raw SQL for now
	query := `
		INSERT INTO registration_tokens (id, token, email, org_id, role, company_profile, expires_at, used, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	r.logger.Info("📝 Executing SQL query", zap.String("query", query))

	_, err := r.db.ExecContext(ctx, query,
		token.ID,
		token.Token,
		token.Email,
		orgID,
		string(token.Role),
		companyProfile,
		time.Unix(token.ExpiresAt, 0),
		token.Used,
		time.Unix(token.CreatedAt, 0),
	)

	if err != nil {
		r.logger.Error("❌ FAILED to execute SQL query",
			zap.Error(err),
			zap.String("token_id", token.ID.String()),
			zap.String("email", token.Email),
		)
		return fmt.Errorf("failed to create registration token: %w", err)
	}

	r.logger.Info("✅ Registration token created successfully in database",
		zap.String("token_id", token.ID.String()),
		zap.String("email", token.Email),
	)
	return nil
}

func (r *RegistrationTokenRepository) GetByToken(tokenStr string) (*domain.RegistrationToken, error) {
	ctx := context.Background()

	query := `
		SELECT id, token, email, org_id, role, company_profile, expires_at, used, created_at
		FROM registration_tokens
		WHERE token = $1
	`

	var token domain.RegistrationToken
	var orgID pgtype.UUID
	var companyProfile *string
	var expiresAt, createdAt time.Time

	err := r.db.QueryRowContext(ctx, query, tokenStr).Scan(
		&token.ID,
		&token.Token,
		&token.Email,
		&orgID,
		&token.Role,
		&companyProfile,
		&expiresAt,
		&token.Used,
		&createdAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		r.logger.Error("Failed to get registration token", zap.Error(err))
		return nil, fmt.Errorf("failed to get registration token: %w", err)
	}

	if orgID.Valid {
		orgUUID := uuid.UUID(orgID.Bytes)
		token.OrgID = &orgUUID
	}

	if companyProfile != nil && *companyProfile != "" {
		profile := domain.CompanyProfile(*companyProfile)
		token.CompanyProfile = &profile
	}

	token.ExpiresAt = expiresAt.Unix()
	token.CreatedAt = createdAt.Unix()

	return &token, nil
}

func (r *RegistrationTokenRepository) MarkAsUsed(tokenStr string) error {
	ctx := context.Background()

	query := `
		UPDATE registration_tokens
		SET used = true
		WHERE token = $1
	`

	result, err := r.db.ExecContext(ctx, query, tokenStr)
	if err != nil {
		r.logger.Error("Failed to mark registration token as used", zap.Error(err))
		return fmt.Errorf("failed to mark registration token as used: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

func (r *RegistrationTokenRepository) ListAll() ([]*domain.RegistrationToken, error) {
	ctx := context.Background()

	query := `
		SELECT id, token, email, org_id, role, company_profile, expires_at, used, created_at
		FROM registration_tokens
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		r.logger.Error("Failed to list registration tokens", zap.Error(err))
		return nil, fmt.Errorf("failed to list registration tokens: %w", err)
	}
	defer rows.Close()

	var tokens []*domain.RegistrationToken
	for rows.Next() {
		var token domain.RegistrationToken
		var orgID pgtype.UUID
		var companyProfile *string
		var expiresAt, createdAt time.Time

		err := rows.Scan(
			&token.ID,
			&token.Token,
			&token.Email,
			&orgID,
			&token.Role,
			&companyProfile,
			&expiresAt,
			&token.Used,
			&createdAt,
		)
		if err != nil {
			r.logger.Error("Failed to scan registration token", zap.Error(err))
			return nil, fmt.Errorf("failed to scan registration token: %w", err)
		}

		if orgID.Valid {
			orgUUID := uuid.UUID(orgID.Bytes)
			token.OrgID = &orgUUID
		}

		if companyProfile != nil && *companyProfile != "" {
			profile := domain.CompanyProfile(*companyProfile)
			token.CompanyProfile = &profile
		}

		token.ExpiresAt = expiresAt.Unix()
		token.CreatedAt = createdAt.Unix()

		tokens = append(tokens, &token)
	}

	return tokens, nil
}

func (r *RegistrationTokenRepository) DeleteExpired() error {
	ctx := context.Background()

	query := `
		DELETE FROM registration_tokens
		WHERE expires_at < NOW()
	`

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		r.logger.Error("Failed to delete expired registration tokens", zap.Error(err))
		return fmt.Errorf("failed to delete expired registration tokens: %w", err)
	}

	return nil
}
