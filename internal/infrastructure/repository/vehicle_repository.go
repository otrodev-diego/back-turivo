package repository

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type VehicleRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewVehicleRepository(db *sql.DB, logger *zap.Logger) *VehicleRepository {
	return &VehicleRepository{
		db:     db,
		logger: logger,
	}
}

func (r *VehicleRepository) Create(vehicle *domain.Vehicle) error {
	ctx := context.Background()

	query := `
		INSERT INTO vehicles (
			id, driver_id, type, brand, model, year, plate, vin, color, capacity,
			insurance_policy, insurance_expires_at, inspection_expires_at, status,
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, NOW(), NOW())
	`

	vehicle.ID = uuid.New().String()
	_, err := r.db.ExecContext(ctx, query,
		vehicle.ID,
		vehicle.DriverID,
		string(vehicle.Type),
		vehicle.Brand,
		vehicle.Model,
		vehicle.Year,
		vehicle.Plate,
		vehicle.VIN,
		vehicle.Color,
		vehicle.Capacity,
		vehicle.InsurancePolicy,
		vehicle.InsuranceExpiresAt,
		vehicle.InspectionExpiresAt,
		string(vehicle.Status),
	)

	if err != nil {
		r.logger.Error("Failed to create vehicle", zap.Error(err))
		return fmt.Errorf("failed to create vehicle: %w", err)
	}

	r.logger.Info("Vehicle created successfully", zap.String("id", vehicle.ID))
	return nil
}

func (r *VehicleRepository) GetByID(id string) (*domain.Vehicle, error) {
	ctx := context.Background()

	query := `
		SELECT 
			v.id, v.driver_id, v.type, v.brand, v.model, v.year, v.plate, v.vin, v.color,
			v.capacity, v.insurance_policy, v.insurance_expires_at, v.inspection_expires_at,
			v.status, v.created_at, v.updated_at,
			d.id, d.first_name, d.last_name, d.phone
		FROM vehicles v
		LEFT JOIN drivers d ON v.driver_id = d.id
		WHERE v.id = $1
	`

	var vehicle domain.Vehicle
	var driverID sql.NullString
	var driverFirstName sql.NullString
	var driverLastName sql.NullString
	var driverPhone sql.NullString

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&vehicle.ID,
		&vehicle.DriverID,
		&vehicle.Type,
		&vehicle.Brand,
		&vehicle.Model,
		&vehicle.Year,
		&vehicle.Plate,
		&vehicle.VIN,
		&vehicle.Color,
		&vehicle.Capacity,
		&vehicle.InsurancePolicy,
		&vehicle.InsuranceExpiresAt,
		&vehicle.InspectionExpiresAt,
		&vehicle.Status,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
		&driverID,
		&driverFirstName,
		&driverLastName,
		&driverPhone,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		r.logger.Error("Failed to get vehicle by ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get vehicle by ID: %w", err)
	}

	// Populate driver info if exists
	if driverID.Valid {
		phone := ""
		if driverPhone.Valid {
			phone = driverPhone.String
		}
		vehicle.Driver = &domain.DriverBasicInfo{
			ID:        driverID.String,
			FirstName: driverFirstName.String,
			LastName:  driverLastName.String,
			Phone:     phone,
		}
	}

	// Get photos
	photos, err := r.GetPhotos(id)
	if err == nil {
		vehicle.Photos = photos
	}

	return &vehicle, nil
}

func (r *VehicleRepository) GetByDriverID(driverID string) (*domain.Vehicle, error) {
	ctx := context.Background()

	query := `
		SELECT 
			id, driver_id, type, brand, model, year, plate, vin, color, capacity,
			insurance_policy, insurance_expires_at, inspection_expires_at, status,
			created_at, updated_at
		FROM vehicles
		WHERE driver_id = $1
		LIMIT 1
	`

	var vehicle domain.Vehicle
	err := r.db.QueryRowContext(ctx, query, driverID).Scan(
		&vehicle.ID,
		&vehicle.DriverID,
		&vehicle.Type,
		&vehicle.Brand,
		&vehicle.Model,
		&vehicle.Year,
		&vehicle.Plate,
		&vehicle.VIN,
		&vehicle.Color,
		&vehicle.Capacity,
		&vehicle.InsurancePolicy,
		&vehicle.InsuranceExpiresAt,
		&vehicle.InspectionExpiresAt,
		&vehicle.Status,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		r.logger.Error("Failed to get vehicle by driver ID", zap.Error(err))
		return nil, fmt.Errorf("failed to get vehicle by driver ID: %w", err)
	}

	return &vehicle, nil
}

func (r *VehicleRepository) List(req domain.ListVehiclesRequest) ([]*domain.Vehicle, int, error) {
	ctx := context.Background()
	r.logger.Info("🚗 === VehicleRepository.List Started ===")

	// Build query with filters
	baseQuery := `
		SELECT 
			v.id, v.driver_id, v.type, v.brand, v.model, v.year, v.plate, v.vin, v.color,
			v.capacity, v.insurance_policy, v.insurance_expires_at, v.inspection_expires_at,
			v.status, v.created_at, v.updated_at,
			d.id, d.first_name, d.last_name, d.phone
		FROM vehicles v
		LEFT JOIN drivers d ON v.driver_id = d.id
	`
	whereClause := ""
	args := []interface{}{}
	argIndex := 1

	// Add search query filter
	if req.Query != nil && *req.Query != "" {
		whereClause = " WHERE (v.brand ILIKE $" + fmt.Sprint(argIndex) + " OR v.model ILIKE $" + fmt.Sprint(argIndex) + " OR v.plate ILIKE $" + fmt.Sprint(argIndex) + ")"
		searchTerm := "%" + *req.Query + "%"
		args = append(args, searchTerm)
		argIndex++
		r.logger.Info("🔍 Filtering by search query", zap.String("query", *req.Query))
	}

	// Add type filter
	if req.Type != nil {
		if whereClause == "" {
			whereClause = " WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += "v.type = $" + fmt.Sprint(argIndex)
		args = append(args, string(*req.Type))
		argIndex++
		r.logger.Info("🔍 Filtering by type", zap.String("type", string(*req.Type)))
	}

	// Add status filter
	if req.Status != nil {
		if whereClause == "" {
			whereClause = " WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += "v.status = $" + fmt.Sprint(argIndex)
		args = append(args, string(*req.Status))
		argIndex++
		r.logger.Info("🔍 Filtering by status", zap.String("status", string(*req.Status)))
	}

	// Add driver filter
	if req.DriverID != nil {
		if whereClause == "" {
			whereClause = " WHERE "
		} else {
			whereClause += " AND "
		}
		whereClause += "v.driver_id = $" + fmt.Sprint(argIndex)
		args = append(args, *req.DriverID)
		argIndex++
		r.logger.Info("🔍 Filtering by driver", zap.String("driver_id", *req.DriverID))
	}

	// Add ORDER BY clause
	orderBy := " ORDER BY v.created_at DESC"
	if req.Sort != "" {
		switch req.Sort {
		case "brand":
			orderBy = " ORDER BY v.brand ASC"
		case "model":
			orderBy = " ORDER BY v.model ASC"
		case "year":
			orderBy = " ORDER BY v.year DESC"
		case "created_at":
			orderBy = " ORDER BY v.created_at DESC"
		}
	}

	// Add pagination
	limitOffset := fmt.Sprintf(" LIMIT %d OFFSET %d", req.PageSize, (req.Page-1)*req.PageSize)

	query := baseQuery + whereClause + orderBy + limitOffset
	r.logger.Info("📝 Executing vehicles query",
		zap.String("query", query),
		zap.Any("args", args),
		zap.Int("page", req.Page),
		zap.Int("page_size", req.PageSize))

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		r.logger.Error("❌ Failed to query vehicles", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list vehicles: %w", err)
	}
	defer rows.Close()

	vehicles := make([]*domain.Vehicle, 0)
	for rows.Next() {
		var vehicle domain.Vehicle
		var driverID sql.NullString
		var driverFirstName sql.NullString
		var driverLastName sql.NullString
		var driverPhone sql.NullString

		err := rows.Scan(
			&vehicle.ID,
			&vehicle.DriverID,
			&vehicle.Type,
			&vehicle.Brand,
			&vehicle.Model,
			&vehicle.Year,
			&vehicle.Plate,
			&vehicle.VIN,
			&vehicle.Color,
			&vehicle.Capacity,
			&vehicle.InsurancePolicy,
			&vehicle.InsuranceExpiresAt,
			&vehicle.InspectionExpiresAt,
			&vehicle.Status,
			&vehicle.CreatedAt,
			&vehicle.UpdatedAt,
			&driverID,
			&driverFirstName,
			&driverLastName,
			&driverPhone,
		)
		if err != nil {
			r.logger.Error("❌ Failed to scan vehicle", zap.Error(err))
			return nil, 0, fmt.Errorf("failed to scan vehicle: %w", err)
		}

		// Populate driver info if exists
		if driverID.Valid {
			phone := ""
			if driverPhone.Valid {
				phone = driverPhone.String
			}
			vehicle.Driver = &domain.DriverBasicInfo{
				ID:        driverID.String,
				FirstName: driverFirstName.String,
				LastName:  driverLastName.String,
				Phone:     phone,
			}
		}

		r.logger.Info("✅ Vehicle scanned",
			zap.String("id", vehicle.ID),
			zap.String("brand", vehicle.Brand),
			zap.String("model", vehicle.Model))
		vehicles = append(vehicles, &vehicle)
	}

	// Get total count for pagination
	countQuery := "SELECT COUNT(*) FROM vehicles v" + whereClause
	var total int
	err = r.db.QueryRowContext(ctx, countQuery, args...).Scan(&total)
	if err != nil {
		r.logger.Error("❌ Failed to count vehicles", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to count vehicles: %w", err)
	}

	r.logger.Info("✅ Vehicles loaded successfully",
		zap.Int("total", total),
		zap.Int("returned", len(vehicles)))

	return vehicles, total, nil
}

func (r *VehicleRepository) Update(id string, req domain.UpdateVehicleRequest) (*domain.Vehicle, error) {
	ctx := context.Background()

	query := `
		UPDATE vehicles
		SET type = COALESCE($2, type),
		    brand = COALESCE($3, brand),
		    model = COALESCE($4, model),
		    year = COALESCE($5, year),
		    plate = COALESCE($6, plate),
		    vin = COALESCE($7, vin),
		    color = COALESCE($8, color),
		    capacity = COALESCE($9, capacity),
		    insurance_policy = COALESCE($10, insurance_policy),
		    insurance_expires_at = COALESCE($11, insurance_expires_at),
		    inspection_expires_at = COALESCE($12, inspection_expires_at),
		    status = COALESCE($13, status),
		    updated_at = NOW()
		WHERE id = $1
		RETURNING id, driver_id, type, brand, model, year, plate, vin, color, capacity,
		          insurance_policy, insurance_expires_at, inspection_expires_at, status,
		          created_at, updated_at
	`

	var vehicle domain.Vehicle
	err := r.db.QueryRowContext(ctx, query,
		id,
		(*string)(req.Type),
		req.Brand,
		req.Model,
		req.Year,
		req.Plate,
		req.VIN,
		req.Color,
		req.Capacity,
		req.InsurancePolicy,
		req.InsuranceExpiresAt,
		req.InspectionExpiresAt,
		(*string)(req.Status),
	).Scan(
		&vehicle.ID,
		&vehicle.DriverID,
		&vehicle.Type,
		&vehicle.Brand,
		&vehicle.Model,
		&vehicle.Year,
		&vehicle.Plate,
		&vehicle.VIN,
		&vehicle.Color,
		&vehicle.Capacity,
		&vehicle.InsurancePolicy,
		&vehicle.InsuranceExpiresAt,
		&vehicle.InspectionExpiresAt,
		&vehicle.Status,
		&vehicle.CreatedAt,
		&vehicle.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, domain.ErrNotFound
		}
		r.logger.Error("Failed to update vehicle", zap.Error(err))
		return nil, fmt.Errorf("failed to update vehicle: %w", err)
	}

	r.logger.Info("Vehicle updated successfully",
		zap.String("id", vehicle.ID))

	return &vehicle, nil
}

func (r *VehicleRepository) Delete(id string) error {
	ctx := context.Background()

	query := `DELETE FROM vehicles WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		r.logger.Error("Failed to delete vehicle", zap.Error(err))
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		r.logger.Error("Failed to get rows affected", zap.Error(err))
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	r.logger.Info("Vehicle deleted successfully", zap.String("id", id))

	return nil
}

func (r *VehicleRepository) AssignToDriver(vehicleID string, driverID *string) error {
	ctx := context.Background()

	// If assigning to a driver, first unassign any other vehicles from that driver
	if driverID != nil {
		_, err := r.db.ExecContext(ctx, `
			UPDATE vehicles 
			SET driver_id = NULL, status = 'AVAILABLE', updated_at = NOW()
			WHERE driver_id = $1 AND id != $2
		`, *driverID, vehicleID)
		if err != nil {
			r.logger.Error("Failed to unassign previous vehicles", zap.Error(err))
			return fmt.Errorf("failed to unassign previous vehicles: %w", err)
		}
	}

	// Now assign (or unassign) the vehicle
	query := `
		UPDATE vehicles
		SET driver_id = $2, updated_at = NOW()
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query, vehicleID, driverID)
	if err != nil {
		r.logger.Error("Failed to assign vehicle to driver", zap.Error(err))
		return fmt.Errorf("failed to assign vehicle to driver: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return domain.ErrNotFound
	}

	if driverID != nil {
		r.logger.Info("Vehicle assigned to driver",
			zap.String("vehicle_id", vehicleID),
			zap.String("driver_id", *driverID))
	} else {
		r.logger.Info("Vehicle unassigned from driver",
			zap.String("vehicle_id", vehicleID))
	}

	return nil
}

func (r *VehicleRepository) AddPhoto(vehicleID string, photoURL string) error {
	ctx := context.Background()

	query := `
		INSERT INTO vehicle_photos (vehicle_id, url, created_at)
		VALUES ($1, $2, NOW())
	`

	_, err := r.db.ExecContext(ctx, query, vehicleID, photoURL)
	if err != nil {
		r.logger.Error("Failed to add vehicle photo", zap.Error(err))
		return fmt.Errorf("failed to add vehicle photo: %w", err)
	}

	return nil
}

func (r *VehicleRepository) RemovePhoto(photoID string) error {
	ctx := context.Background()

	query := `DELETE FROM vehicle_photos WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, photoID)
	if err != nil {
		r.logger.Error("Failed to remove vehicle photo", zap.Error(err))
		return fmt.Errorf("failed to remove vehicle photo: %w", err)
	}

	return nil
}

func (r *VehicleRepository) GetPhotos(vehicleID string) ([]string, error) {
	ctx := context.Background()

	query := `SELECT url FROM vehicle_photos WHERE vehicle_id = $1 ORDER BY created_at`

	rows, err := r.db.QueryContext(ctx, query, vehicleID)
	if err != nil {
		r.logger.Error("Failed to get vehicle photos", zap.Error(err))
		return nil, fmt.Errorf("failed to get vehicle photos: %w", err)
	}
	defer rows.Close()

	var photos []string
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, fmt.Errorf("failed to scan photo URL: %w", err)
		}
		photos = append(photos, url)
	}

	return photos, nil
}

