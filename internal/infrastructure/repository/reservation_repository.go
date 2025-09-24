package repository

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"turivo-backend/internal/domain"
	"turivo-backend/internal/infrastructure/db/sqlc"
)

type ReservationRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewReservationRepository(db *pgxpool.Pool) *ReservationRepository {
	return &ReservationRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *ReservationRepository) Create(reservation *domain.Reservation) error {
	ctx := context.Background()

	var amount pgtype.Numeric
	if reservation.Amount != nil {
		// Convert float64 to pgtype.Numeric properly
		amount = pgtype.Numeric{
			Int:   big.NewInt(int64(*reservation.Amount * 100)), // Convert to cents to preserve precision
			Exp:   -2,                                           // Two decimal places
			Valid: true,
		}
	}

	var userID pgtype.UUID
	if reservation.UserID != nil {
		userID = pgtype.UUID{Bytes: *reservation.UserID, Valid: true}
	}

	var orgID pgtype.UUID
	if reservation.OrgID != nil {
		orgID = pgtype.UUID{Bytes: *reservation.OrgID, Valid: true}
	}

	dbReservation, err := r.queries.CreateReservation(ctx, sqlc.CreateReservationParams{
		ID:          reservation.ID,
		UserID:      userID,
		OrgID:       orgID,
		Pickup:      reservation.Pickup,
		Destination: reservation.Destination,
		Datetime:    pgtype.Timestamptz{Time: reservation.DateTime, Valid: true},
		Passengers:  int32(reservation.Passengers),
		Status:      sqlc.ReservationStatus(reservation.Status),
		Amount:      amount,
		Notes:       reservation.Notes,
	})
	if err != nil {
		return fmt.Errorf("failed to create reservation: %w", err)
	}

	// Update reservation with timestamps
	reservation.CreatedAt = dbReservation.CreatedAt.Time
	reservation.UpdatedAt = dbReservation.UpdatedAt.Time

	return nil
}

func (r *ReservationRepository) GetByID(id string) (*domain.Reservation, error) {
	ctx := context.Background()

	dbReservation, err := r.queries.GetReservationByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrReservationNotFound
		}
		return nil, fmt.Errorf("failed to get reservation by ID: %w", err)
	}

	return r.mapToDomainReservation(dbReservation), nil
}

func (r *ReservationRepository) List(req domain.ListReservationsRequest) ([]*domain.Reservation, int, error) {
	ctx := context.Background()

	// For now, use a simple approach: if no status filter, use a valid status that won't match anything
	// This is a temporary workaround until we fix the SQL query to handle NULL properly

	var statusFilter sqlc.ReservationStatus = "NONEXISTENT" // Default value that won't match any real status
	if req.Status != nil {
		statusFilter = sqlc.ReservationStatus(*req.Status)
	}

	// Convert search query
	var searchQuery string
	if req.Query != nil {
		searchQuery = *req.Query
	}

	// Convert user_id
	var userID pgtype.UUID
	if req.UserID != nil {
		userID = pgtype.UUID{Bytes: *req.UserID, Valid: true}
	}

	// Convert org_id (empty for now, but could be used later)
	var orgID pgtype.UUID

	// Convert sort
	sort := req.Sort
	if sort == "" {
		sort = "datetime"
	}

	// Calculate offset
	offset := (req.Page - 1) * req.PageSize

	// If no status filter is provided, we need to get all statuses
	// Let's use a different approach: build a custom query or use multiple calls
	var dbReservations []sqlc.Reservation
	var err error

	if req.Status == nil {
		// No status filter - we need all reservations for this user
		// Since the SQL query doesn't handle NULL status properly, we'll call it multiple times
		allStatuses := []sqlc.ReservationStatus{
			sqlc.ReservationStatusACTIVA,
			sqlc.ReservationStatusPROGRAMADA,
			sqlc.ReservationStatusCOMPLETADA,
			sqlc.ReservationStatusCANCELADA,
		}

		var allReservations []sqlc.Reservation
		for _, status := range allStatuses {
			reservations, err := r.queries.ListReservations(ctx, sqlc.ListReservationsParams{
				Column1: searchQuery,
				Column2: status,
				Column3: userID,
				Column4: orgID,
				Column5: sort,
				Limit:   int32(req.PageSize * 2), // Get more to account for multiple statuses
				Offset:  0,                       // Reset offset for each query
			})
			if err != nil {
				return nil, 0, fmt.Errorf("failed to list reservations: %w", err)
			}
			allReservations = append(allReservations, reservations...)
		}

		// Sort the combined results by datetime desc
		// For now, just use the first pageSize results
		if len(allReservations) > req.PageSize {
			dbReservations = allReservations[:req.PageSize]
		} else {
			dbReservations = allReservations
		}
	} else {
		// Status filter provided
		dbReservations, err = r.queries.ListReservations(ctx, sqlc.ListReservationsParams{
			Column1: searchQuery,
			Column2: statusFilter,
			Column3: userID,
			Column4: orgID,
			Column5: sort,
			Limit:   int32(req.PageSize),
			Offset:  int32(offset),
		})
		if err != nil {
			return nil, 0, fmt.Errorf("failed to list reservations: %w", err)
		}
	}

	// Convert to domain reservations
	reservations := make([]*domain.Reservation, len(dbReservations))
	for i, dbRes := range dbReservations {
		reservation := r.mapToDomainReservation(dbRes)
		reservations[i] = reservation
	}

	// TODO: Get total count - for now estimate based on results
	total := len(reservations)
	if len(reservations) == req.PageSize {
		// If we got a full page, there might be more
		total = req.Page*req.PageSize + 1
	}

	return reservations, total, nil
}

func (r *ReservationRepository) Update(id string, req domain.UpdateReservationRequest) (*domain.Reservation, error) {
	// TODO: Implement when UpdateReservation query is properly generated by SQLC
	return nil, fmt.Errorf("update reservation not yet implemented")
}

func (r *ReservationRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteReservation(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %w", err)
	}

	return nil
}

func (r *ReservationRepository) ChangeStatus(id string, newStatus domain.ReservationStatus) error {
	ctx := context.Background()

	_, err := r.queries.UpdateReservationStatus(ctx, sqlc.UpdateReservationStatusParams{
		ID:     id,
		Status: sqlc.ReservationStatus(newStatus),
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrReservationNotFound
		}
		return fmt.Errorf("failed to change reservation status: %w", err)
	}

	return nil
}

func (r *ReservationRepository) GetTimeline(id string) ([]domain.TimelineEvent, error) {
	// TODO: Implement when timeline queries are properly generated by SQLC
	return []domain.TimelineEvent{}, nil
}

func (r *ReservationRepository) AddTimelineEvent(id string, event domain.TimelineEvent) error {
	// TODO: Implement when timeline queries are properly generated by SQLC
	return fmt.Errorf("add timeline event not yet implemented")
}

func (r *ReservationRepository) GenerateID() string {
	// Simple implementation - in production, you might want a more sophisticated approach
	timestamp := time.Now().Unix()
	return fmt.Sprintf("RSV-%d", timestamp%100000)
}

func (r *ReservationRepository) mapToDomainReservation(dbReservation sqlc.Reservation) *domain.Reservation {
	reservation := &domain.Reservation{
		ID:          dbReservation.ID,
		Pickup:      dbReservation.Pickup,
		Destination: dbReservation.Destination,
		Passengers:  int(dbReservation.Passengers),
		Status:      domain.ReservationStatus(dbReservation.Status),
		Notes:       dbReservation.Notes,
		CreatedAt:   dbReservation.CreatedAt.Time,
		UpdatedAt:   dbReservation.UpdatedAt.Time,
	}

	// Convert pgtype.UUID to *uuid.UUID
	if dbReservation.UserID.Valid {
		var userID uuid.UUID
		copy(userID[:], dbReservation.UserID.Bytes[:])
		reservation.UserID = &userID
	}

	if dbReservation.OrgID.Valid {
		var orgID uuid.UUID
		copy(orgID[:], dbReservation.OrgID.Bytes[:])
		reservation.OrgID = &orgID
	}

	if dbReservation.Datetime.Valid {
		reservation.DateTime = dbReservation.Datetime.Time
	}

	if dbReservation.Amount.Valid {
		// TODO: Implement proper numeric conversion
		amount := 0.0
		reservation.Amount = &amount
	}

	// Set assigned driver ID if present
	if dbReservation.AssignedDriverID != nil {
		reservation.AssignedDriverID = dbReservation.AssignedDriverID
	}

	return reservation
}

// AssignDriver assigns a driver to a reservation
func (r *ReservationRepository) AssignDriver(reservationID string, driverID string) error {
	ctx := context.Background()

	// For now, we'll use a simple update query
	// In production, you might want to validate the driver exists first
	query := `
		UPDATE reservations 
		SET assigned_driver_id = $1, updated_at = NOW()
		WHERE id = $2`

	result, err := r.db.Exec(ctx, query, driverID, reservationID)
	if err != nil {
		return fmt.Errorf("failed to assign driver: %w", err)
	}

	if result.RowsAffected() == 0 {
		return domain.ErrReservationNotFound
	}

	return nil
}
