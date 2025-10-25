package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"turivo-backend/internal/domain"
)

type DriverDashboardRepository struct {
	db *pgxpool.Pool
}

func NewDriverDashboardRepository(db *pgxpool.Pool) *DriverDashboardRepository {
	return &DriverDashboardRepository{
		db: db,
	}
}

// GetByUserID gets a driver by user ID
func (r *DriverDashboardRepository) GetByUserID(userID string) (*domain.Driver, error) {
	ctx := context.Background()

	query := `
		SELECT d.id, d.first_name, d.last_name, d.rut_or_dni, d.birth_date, 
		       d.phone, d.email, d.photo_url, d.status, d.created_at, d.updated_at,
		       d.user_id, d.company_id, d.vehicle_id
		FROM drivers d
		WHERE d.user_id = $1
	`

	row := r.db.QueryRow(ctx, query, userID)

	var driver domain.Driver
	var birthDate pgtype.Timestamp
	var phone, email, photoURL pgtype.Text
	var userIDPtr, companyIDPtr, vehicleIDPtr pgtype.Text

	err := row.Scan(
		&driver.ID,
		&driver.FirstName,
		&driver.LastName,
		&driver.RutOrDNI,
		&birthDate,
		&phone,
		&email,
		&photoURL,
		&driver.Status,
		&driver.CreatedAt,
		&driver.UpdatedAt,
		&userIDPtr,
		&companyIDPtr,
		&vehicleIDPtr,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get driver by user ID: %w", err)
	}

	// Map optional fields
	if birthDate.Valid {
		driver.BirthDate = &birthDate.Time
	}
	if phone.Valid {
		driver.Phone = &phone.String
	}
	if email.Valid {
		driver.Email = &email.String
	}
	if photoURL.Valid {
		driver.PhotoURL = &photoURL.String
	}
	if userIDPtr.Valid {
		driver.UserID = &userIDPtr.String
	}
	if companyIDPtr.Valid {
		driver.CompanyID = &companyIDPtr.String
	}
	if vehicleIDPtr.Valid {
		driver.VehicleID = &vehicleIDPtr.String
	}

	return &driver, nil
}

// GetDriverTrips gets trips for a specific driver
func (r *DriverDashboardRepository) GetDriverTrips(driverID string) ([]*domain.Reservation, error) {
	ctx := context.Background()

	query := `
		SELECT id, user_id, pickup, destination, datetime, passengers, 
		       status, distance_km, fare, notes, assigned_driver_id, 
		       created_at, updated_at
		FROM reservations 
		WHERE assigned_driver_id = $1
		ORDER BY datetime DESC
	`

	rows, err := r.db.Query(ctx, query, driverID)
	if err != nil {
		return nil, fmt.Errorf("failed to get driver trips: %w", err)
	}
	defer rows.Close()

	var trips []*domain.Reservation
	for rows.Next() {
		var trip domain.Reservation
		var distanceKm, fare pgtype.Numeric
		var notes pgtype.Text
		var assignedDriverID pgtype.Text

		err := rows.Scan(
			&trip.ID,
			&trip.UserID,
			&trip.Pickup,
			&trip.Destination,
			&trip.DateTime,
			&trip.Passengers,
			&trip.Status,
			&distanceKm,
			&fare,
			&notes,
			&assignedDriverID,
			&trip.CreatedAt,
			&trip.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan trip: %w", err)
		}

		// Map optional fields
		if distanceKm.Valid {
			amount := float64(distanceKm.Int.Int64()) / 100.0 // Convert from cents to dollars
			trip.Amount = &amount
		}
		if fare.Valid {
			amount := float64(fare.Int.Int64()) / 100.0 // Convert from cents to dollars
			trip.Amount = &amount
		}
		if notes.Valid {
			trip.Notes = &notes.String
		}
		if assignedDriverID.Valid {
			trip.AssignedDriverID = &assignedDriverID.String
		}

		trips = append(trips, &trip)
	}

	return trips, nil
}

// GetDriverVehicle gets the vehicle assigned to a driver
func (r *DriverDashboardRepository) GetDriverVehicle(driverID string) (*domain.Vehicle, error) {
	ctx := context.Background()

	query := `
		SELECT v.id, v.make, v.model, v.year, v.license_plate, v.color, 
		       v.capacity, v.status, v.created_at, v.updated_at
		FROM vehicles v
		INNER JOIN drivers d ON d.vehicle_id = v.id
		WHERE d.id = $1
	`

	row := r.db.QueryRow(ctx, query, driverID)

	var vehicle domain.Vehicle
	var year pgtype.Int4
	var capacity pgtype.Int4

	err := row.Scan(
		&vehicle.ID,
		&vehicle.Brand,
		&vehicle.Model,
		&year,
		&vehicle.Plate,
		&vehicle.Color,
		&capacity,
		&vehicle.Status,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		return nil, fmt.Errorf("failed to get driver vehicle: %w", err)
	}

	// Map numeric fields
	if year.Valid {
		yearValue := int(year.Int32)
		vehicle.Year = &yearValue
	}
	if capacity.Valid {
		capacityValue := int(capacity.Int32)
		vehicle.Capacity = &capacityValue
	}

	return &vehicle, nil
}

// UpdateTripStatus updates the status of a trip
func (r *DriverDashboardRepository) UpdateTripStatus(driverID, tripID, status string) error {
	ctx := context.Background()

	// First verify that the trip belongs to this driver
	verifyQuery := `
		SELECT id FROM reservations 
		WHERE id = $1 AND assigned_driver_id = $2
	`

	var tripExists string
	err := r.db.QueryRow(ctx, verifyQuery, tripID, driverID).Scan(&tripExists)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.ErrNotFound
		}
		return fmt.Errorf("failed to verify trip ownership: %w", err)
	}

	// Update the trip status
	updateQuery := `
		UPDATE reservations 
		SET status = $1, updated_at = NOW()
		WHERE id = $2 AND assigned_driver_id = $3
	`

	result, err := r.db.Exec(ctx, updateQuery, status, tripID, driverID)
	if err != nil {
		return fmt.Errorf("failed to update trip status: %w", err)
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	return nil
}
