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

type UserRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *UserRepository) Create(user *domain.User) error {
	ctx := context.Background()

	var orgID pgtype.UUID
	if user.OrgID != nil {
		orgID = pgtype.UUID{Bytes: *user.OrgID, Valid: true}
	}

	var companyProfile *string
	if user.CompanyProfile != nil {
		profileStr := string(*user.CompanyProfile)
		companyProfile = &profileStr
	}

	dbUser, err := r.queries.CreateUser(ctx, sqlc.CreateUserParams{
		Name:           user.Name,
		Email:          user.Email,
		PasswordHash:   user.PasswordHash,
		Role:           sqlc.UserRole(user.Role),
		OrgID:          orgID,
		CompanyProfile: companyProfile,
	})
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Update user with generated ID and timestamps
	copy(user.ID[:], dbUser.ID.Bytes[:])
	user.CreatedAt = dbUser.CreatedAt.Time
	user.UpdatedAt = dbUser.UpdatedAt.Time

	return nil
}

func (r *UserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	ctx := context.Background()

	dbUser, err := r.queries.GetUserByID(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	return r.mapToDomainUser(dbUser), nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	ctx := context.Background()

	dbUser, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrUserNotFound
		}
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	return r.mapToDomainUser(dbUser), nil
}

func (r *UserRepository) List(req domain.ListUsersRequest) ([]*domain.User, int, error) {
	ctx := context.Background()

	// Build query with filters
	baseQuery := `
		SELECT id, name, email, password_hash, role, status, org_id, created_at, updated_at
		FROM users
	`
	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	// Filter by organization ID if provided
	if req.OrgID != nil {
		whereClause += " WHERE org_id = $" + fmt.Sprint(argIndex)
		args = append(args, *req.OrgID)
		argIndex++
	}

	// Add search query filter
	if req.Query != nil && *req.Query != "" {
		if whereClause == "" {
			whereClause += " WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += "(name ILIKE $" + fmt.Sprint(argIndex) + " OR email ILIKE $" + fmt.Sprint(argIndex) + ")"
		searchTerm := "%" + *req.Query + "%"
		args = append(args, searchTerm)
		argIndex++
	}

	// Add role filter
	if req.Role != nil {
		if whereClause == "" {
			whereClause += " WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += "role = $" + fmt.Sprint(argIndex)
		args = append(args, string(*req.Role))
		argIndex++
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
	}

	// Add ORDER BY clause
	orderBy := " ORDER BY created_at DESC"
	if req.Sort != "" {
		switch req.Sort {
		case "name":
			orderBy = " ORDER BY name ASC"
		case "email":
			orderBy = " ORDER BY email ASC"
		case "created_at":
			orderBy = " ORDER BY created_at DESC"
		}
	}

	// Add pagination
	limitOffset := fmt.Sprintf(" LIMIT %d OFFSET %d", req.PageSize, (req.Page-1)*req.PageSize)

	query := baseQuery + whereClause + orderBy + limitOffset

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*domain.User
	for rows.Next() {
		var userID, orgID pgtype.UUID
		var name, email, passwordHash, role, status string
		var createdAt, updatedAt pgtype.Timestamptz

		err := rows.Scan(&userID, &name, &email, &passwordHash, &role, &status, &orgID, &createdAt, &updatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}

		user := &domain.User{
			Name:         name,
			Email:        email,
			PasswordHash: passwordHash,
			Role:         domain.UserRole(role),
			Status:       domain.UserStatus(status),
		}

		// Convert UUIDs
		if userID.Valid {
			var id uuid.UUID
			copy(id[:], userID.Bytes[:])
			user.ID = id
		}

		if orgID.Valid {
			var id uuid.UUID
			copy(id[:], orgID.Bytes[:])
			user.OrgID = &id
		}

		if createdAt.Valid {
			user.CreatedAt = createdAt.Time
		}

		if updatedAt.Valid {
			user.UpdatedAt = updatedAt.Time
		}

		users = append(users, user)
	}

	// Get total count for pagination
	countQuery := "SELECT COUNT(*) FROM users" + whereClause
	var total int
	err = r.db.QueryRow(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	return users, total, nil
}

func (r *UserRepository) Update(id uuid.UUID, req domain.UpdateUserRequest) (*domain.User, error) {
	ctx := context.Background()

	// Convert domain UUID to pgtype.UUID
	var pgID pgtype.UUID
	copy(pgID.Bytes[:], id[:])
	pgID.Valid = true

	// Get the current user to use as base values
	currentUser, err := r.queries.GetUserByID(ctx, pgID)
	if err != nil {
		return nil, fmt.Errorf("failed to get current user: %w", err)
	}

	// Prepare parameters for the update, using current values as defaults
	params := sqlc.UpdateUserParams{
		ID:             pgID,
		Name:           currentUser.Name,
		Email:          currentUser.Email,
		PasswordHash:   currentUser.PasswordHash,
		Role:           currentUser.Role,
		Status:         currentUser.Status,
		OrgID:          currentUser.OrgID,
		CompanyProfile: currentUser.CompanyProfile,
	}

	// Override only the fields that are provided (not nil)
	if req.Name != nil {
		params.Name = *req.Name
	}
	if req.Email != nil {
		params.Email = *req.Email
	}
	if req.Password != nil {
		params.PasswordHash = *req.Password
	}
	if req.Role != nil {
		params.Role = sqlc.UserRole(*req.Role)
	}
	if req.Status != nil {
		params.Status = sqlc.UserStatus(*req.Status)
	}
	if req.OrgID != nil {
		var pgOrgID pgtype.UUID
		copy(pgOrgID.Bytes[:], (*req.OrgID)[:])
		pgOrgID.Valid = true
		params.OrgID = pgOrgID
	}
	if req.CompanyProfile != nil {
		profileStr := string(*req.CompanyProfile)
		params.CompanyProfile = &profileStr
	}

	// Execute the update
	updatedUser, err := r.queries.UpdateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	// Convert back to domain.User
	user := &domain.User{
		Name:         updatedUser.Name,
		Email:        updatedUser.Email,
		PasswordHash: updatedUser.PasswordHash,
		Role:         domain.UserRole(updatedUser.Role),
		Status:       domain.UserStatus(updatedUser.Status),
	}

	// Convert UUIDs
	if updatedUser.ID.Valid {
		var userID uuid.UUID
		copy(userID[:], updatedUser.ID.Bytes[:])
		user.ID = userID
	}

	if updatedUser.OrgID.Valid {
		var orgID uuid.UUID
		copy(orgID[:], updatedUser.OrgID.Bytes[:])
		user.OrgID = &orgID
	}

	if updatedUser.CreatedAt.Valid {
		user.CreatedAt = updatedUser.CreatedAt.Time
	}

	if updatedUser.UpdatedAt.Valid {
		user.UpdatedAt = updatedUser.UpdatedAt.Time
	}

	if updatedUser.CompanyProfile != nil {
		profile := domain.CompanyProfile(*updatedUser.CompanyProfile)
		user.CompanyProfile = &profile
	}

	return user, nil
}

func (r *UserRepository) UpdateUser(user *domain.User) error {
	ctx := context.Background()

	// Convert domain UUID to pgtype.UUID
	var pgID pgtype.UUID
	copy(pgID.Bytes[:], user.ID[:])
	pgID.Valid = true

	var orgID pgtype.UUID
	if user.OrgID != nil {
		copy(orgID.Bytes[:], (*user.OrgID)[:])
		orgID.Valid = true
	}

	var companyProfile *string
	if user.CompanyProfile != nil {
		profileStr := string(*user.CompanyProfile)
		companyProfile = &profileStr
	}

	// Execute the update
	_, err := r.queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:             pgID,
		Name:           user.Name,
		Email:          user.Email,
		PasswordHash:   user.PasswordHash,
		Role:           sqlc.UserRole(user.Role),
		Status:         sqlc.UserStatus(user.Status),
		OrgID:          orgID,
		CompanyProfile: companyProfile,
	})
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	return nil
}

func (r *UserRepository) Delete(id uuid.UUID) error {
	ctx := context.Background()

	err := r.queries.DeleteUser(ctx, pgtype.UUID{Bytes: id, Valid: true})
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *UserRepository) mapToDomainUser(dbUser sqlc.User) *domain.User {
	var userID uuid.UUID
	copy(userID[:], dbUser.ID.Bytes[:])

	user := &domain.User{
		ID:           userID,
		Name:         dbUser.Name,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		Role:         domain.UserRole(dbUser.Role),
		Status:       domain.UserStatus(dbUser.Status),
		CreatedAt:    dbUser.CreatedAt.Time,
		UpdatedAt:    dbUser.UpdatedAt.Time,
	}

	if dbUser.OrgID.Valid {
		var orgID uuid.UUID
		copy(orgID[:], dbUser.OrgID.Bytes[:])
		user.OrgID = &orgID
	}

	if dbUser.CompanyProfile != nil && *dbUser.CompanyProfile != "" {
		profile := domain.CompanyProfile(*dbUser.CompanyProfile)
		user.CompanyProfile = &profile
	}

	return user
}
