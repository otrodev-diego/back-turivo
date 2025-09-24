package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type CompanyRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewCompanyRepository(db *sql.DB, logger *zap.Logger) *CompanyRepository {
	return &CompanyRepository{
		db:     db,
		logger: logger,
	}
}

func (r *CompanyRepository) Create(company *domain.Company) error {
	ctx := context.Background()

	query := `
		INSERT INTO companies (id, name, rut, contact_email, status, sector, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
	`

	company.ID = uuid.New()
	_, err := r.db.ExecContext(ctx, query,
		company.ID,
		company.Name,
		company.RUT,
		company.ContactEmail,
		string(company.Status),
		string(company.Sector),
	)

	if err != nil {
		r.logger.Error("Failed to create company", zap.Error(err))
		return fmt.Errorf("failed to create company: %w", err)
	}

	return nil
}

func (r *CompanyRepository) GetByID(id uuid.UUID) (*domain.Company, error) {
	ctx := context.Background()

	query := `
		SELECT id, name, rut, contact_email, status, sector, created_at, updated_at
		FROM companies
		WHERE id = $1
	`

	var company domain.Company
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&company.ID,
		&company.Name,
		&company.RUT,
		&company.ContactEmail,
		&company.Status,
		&company.Sector,
		&company.CreatedAt,
		&company.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		r.logger.Error("Failed to get company by ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get company by ID: %w", err)
	}

	return &company, nil
}

func (r *CompanyRepository) GetByRUT(rut string) (*domain.Company, error) {
	ctx := context.Background()

	query := `
		SELECT id, name, rut, contact_email, status, sector, created_at, updated_at
		FROM companies
		WHERE rut = $1
	`

	var company domain.Company
	err := r.db.QueryRowContext(ctx, query, rut).Scan(
		&company.ID,
		&company.Name,
		&company.RUT,
		&company.ContactEmail,
		&company.Status,
		&company.Sector,
		&company.CreatedAt,
		&company.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		r.logger.Error("Failed to get company by RUT", zap.Error(err))
		return nil, fmt.Errorf("failed to get company by RUT: %w", err)
	}

	return &company, nil
}

func (r *CompanyRepository) List(req domain.ListCompaniesRequest) ([]*domain.Company, int, error) {
	ctx := context.Background()
	r.logger.Info("🏢 === CompanyRepository.List Started ===")

	// Build query with filters
	baseQuery := `
		SELECT id, name, rut, contact_email, status, sector, created_at, updated_at
		FROM companies
	`
	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	// Filter by organization ID if provided
	if req.OrgID != nil {
		whereClause += " WHERE id = $" + fmt.Sprint(argIndex)
		args = append(args, *req.OrgID)
		argIndex++
		r.logger.Info("🔍 Filtering by organization ID", zap.String("org_id", req.OrgID.String()))
	}

	// Add search query filter
	if req.Query != nil && *req.Query != "" {
		if whereClause == "" {
			whereClause += " WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += "(name ILIKE $" + fmt.Sprint(argIndex) + " OR rut ILIKE $" + fmt.Sprint(argIndex) + ")"
		searchTerm := "%" + *req.Query + "%"
		args = append(args, searchTerm)
		argIndex++
		r.logger.Info("🔍 Filtering by search query", zap.String("query", *req.Query))
	}

	// Add status filter
	if req.Status != nil {
		if whereClause == "" {
			whereClause += " WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += "status = $" + fmt.Sprint(argIndex)
		args = append(args, string(*req.Status))
		argIndex++
		r.logger.Info("🔍 Filtering by status", zap.String("status", string(*req.Status)))
	}

	// Add sector filter
	if req.Sector != nil {
		if whereClause == "" {
			whereClause += " WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += "sector = $" + fmt.Sprint(argIndex)
		args = append(args, string(*req.Sector))
		argIndex++
		r.logger.Info("🔍 Filtering by sector", zap.String("sector", string(*req.Sector)))
	}

	// Add ORDER BY clause
	orderBy := " ORDER BY created_at DESC"
	if req.Sort != "" {
		switch req.Sort {
		case "name":
			orderBy = " ORDER BY name ASC"
		case "rut":
			orderBy = " ORDER BY rut ASC"
		case "created_at":
			orderBy = " ORDER BY created_at DESC"
		}
	}

	// Add pagination
	limitOffset := fmt.Sprintf(" LIMIT %d OFFSET %d", req.PageSize, (req.Page-1)*req.PageSize)

	query := baseQuery + whereClause + orderBy + limitOffset
	r.logger.Info("📝 Executing companies query",
		zap.String("query", query),
		zap.Any("args", args),
		zap.Int("page", req.Page),
		zap.Int("page_size", req.PageSize))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("❌ Failed to query companies", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list companies: %w", err)
	}
	defer rows.Close()

	var companies []*domain.Company
	for rows.Next() {
		var company domain.Company
		err := rows.Scan(
			&company.ID,
			&company.Name,
			&company.RUT,
			&company.ContactEmail,
			&company.Status,
			&company.Sector,
			&company.CreatedAt,
			&company.UpdatedAt,
		)
		if err != nil {
			r.logger.Error("❌ Failed to scan company", zap.Error(err))
			return nil, 0, fmt.Errorf("failed to scan company: %w", err)
		}
		r.logger.Info("✅ Company scanned",
			zap.String("id", company.ID.String()),
			zap.String("name", company.Name),
		)
		companies = append(companies, &company)
	}

	// Get total count for pagination
	countQuery := "SELECT COUNT(*) FROM companies" + whereClause
	var total int
	err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		r.logger.Error("❌ Failed to count companies", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to count companies: %w", err)
	}

	r.logger.Info("✅ Companies loaded successfully",
		zap.Int("total", total),
		zap.Int("returned", len(companies)))

	return companies, total, nil
}

func (r *CompanyRepository) Update(id uuid.UUID, req domain.UpdateCompanyRequest) (*domain.Company, error) {
	ctx := context.Background()

	query := `
		UPDATE companies
		SET name = COALESCE($2, name),
		    rut = COALESCE($3, rut),
		    contact_email = COALESCE($4, contact_email),
		    status = COALESCE($5, status),
		    sector = COALESCE($6, sector),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, rut, contact_email, status, sector, created_at, updated_at
	`

	var company domain.Company
	err := r.db.QueryRowContext(ctx, query,
		id,
		req.Name,
		req.RUT,
		req.ContactEmail,
		(*string)(req.Status),
		(*string)(req.Sector),
	).Scan(
		&company.ID,
		&company.Name,
		&company.RUT,
		&company.ContactEmail,
		&company.Status,
		&company.Sector,
		&company.CreatedAt,
		&company.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		r.logger.Error("Failed to update company", zap.Error(err))
		return nil, fmt.Errorf("failed to update company: %w", err)
	}

	r.logger.Info("Company updated successfully",
		zap.String("id", company.ID.String()),
		zap.String("name", company.Name))

	return &company, nil
}

func (r *CompanyRepository) Delete(id uuid.UUID) error {
	ctx := context.Background()

	query := `DELETE FROM companies WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete company", zap.Error(err))
		return fmt.Errorf("failed to delete company: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", zap.Error(err))
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	r.logger.Info("Company deleted successfully", zap.String("id", id.String()))

	return nil
}
