package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"turivo-backend/internal/domain"
	"turivo-backend/internal/infrastructure/db/sqlc"
)

type DriverRepository struct {
	db      *pgxpool.Pool
	queries *sqlc.Queries
}

func NewDriverRepository(db *pgxpool.Pool) *DriverRepository {
	return &DriverRepository{
		db:      db,
		queries: sqlc.New(db),
	}
}

func (r *DriverRepository) Create(driver *domain.Driver) error {
	ctx := context.Background()

	var birthDate pgtype.Date
	if driver.BirthDate != nil {
		birthDate = pgtype.Date{Time: *driver.BirthDate, Valid: true}
	}

	dbDriver, err := r.queries.CreateDriver(ctx, sqlc.CreateDriverParams{
		ID:        driver.ID,
		FirstName: driver.FirstName,
		LastName:  driver.LastName,
		RutOrDni:  driver.RutOrDNI,
		BirthDate: birthDate,
		Phone:     driver.Phone,
		Email:     driver.Email,
		PhotoUrl:  driver.PhotoURL,
		Status:    sqlc.DriverStatus(driver.Status),
	})
	if err != nil {
		return fmt.Errorf("failed to create driver: %w", err)
	}

	// Update driver with timestamps
	driver.CreatedAt = dbDriver.CreatedAt.Time
	driver.UpdatedAt = dbDriver.UpdatedAt.Time

	return nil
}

func (r *DriverRepository) GetByID(id string) (*domain.Driver, error) {
	ctx := context.Background()

	dbDriver, err := r.queries.GetDriverByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDriverNotFound
		}
		return nil, fmt.Errorf("failed to get driver by ID: %w", err)
	}

	return r.mapToDomainDriver(dbDriver), nil
}

func (r *DriverRepository) List(req domain.ListDriversRequest) ([]*domain.Driver, int, error) {
	ctx := context.Background()

	offset := (req.Page - 1) * req.PageSize

	var query *string
	if req.Query != nil && *req.Query != "" {
		query = req.Query
	}

	var status *sqlc.DriverStatus
	if req.Status != nil {
		sqlcStatus := sqlc.DriverStatus(*req.Status)
		status = &sqlcStatus
	}

	// Get drivers
	queryParam := ""
	if query != nil {
		queryParam = *query
	}
	
	statusParam := sqlc.DriverStatus("")
	if status != nil {
		statusParam = *status
	}

	dbDrivers, err := r.queries.ListDrivers(ctx, sqlc.ListDriversParams{
		Column1: queryParam,
		Column2: statusParam,
		Column3: req.Sort,
		Limit:   int32(req.PageSize),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list drivers: %w", err)
	}

	// Get total count
	total, err := r.queries.CountDrivers(ctx, sqlc.CountDriversParams{
		Column1: queryParam,
		Column2: statusParam,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count drivers: %w", err)
	}

	// Map to domain drivers
	drivers := make([]*domain.Driver, len(dbDrivers))
	for i, dbDriver := range dbDrivers {
		drivers[i] = r.mapToDomainDriverFromList(dbDriver)
	}

	return drivers, int(total), nil
}

func (r *DriverRepository) Update(id string, req domain.UpdateDriverRequest) (*domain.Driver, error) {
	ctx := context.Background()

	// Variables removed as they are not needed - we use the values directly in the update logic

	// Get current driver to use existing values for nil fields
	currentDriver, err := r.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Use current values if request fields are nil
	firstName := currentDriver.FirstName
	if req.FirstName != nil {
		firstName = *req.FirstName
	}

	lastName := currentDriver.LastName
	if req.LastName != nil {
		lastName = *req.LastName
	}

	rutOrDNI := currentDriver.RutOrDNI
	if req.RutOrDNI != nil {
		rutOrDNI = *req.RutOrDNI
	}

	var updateBirthDate pgtype.Date
	if req.BirthDate != nil {
		updateBirthDate = pgtype.Date{Time: *req.BirthDate, Valid: true}
	} else if currentDriver.BirthDate != nil {
		updateBirthDate = pgtype.Date{Time: *currentDriver.BirthDate, Valid: true}
	}

	phone := currentDriver.Phone
	if req.Phone != nil {
		phone = req.Phone
	}

	email := currentDriver.Email
	if req.Email != nil {
		email = req.Email
	}

	photoURL := currentDriver.PhotoURL
	if req.PhotoURL != nil {
		photoURL = req.PhotoURL
	}

	driverStatus := sqlc.DriverStatus(currentDriver.Status)
	if req.Status != nil {
		driverStatus = sqlc.DriverStatus(*req.Status)
	}

	_, err = r.queries.UpdateDriver(ctx, sqlc.UpdateDriverParams{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		RutOrDni:  rutOrDNI,
		BirthDate: updateBirthDate,
		Phone:     phone,
		Email:     email,
		PhotoUrl:  photoURL,
		Status:    driverStatus,
	})
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, domain.ErrDriverNotFound
		}
		return nil, fmt.Errorf("failed to update driver: %w", err)
	}

	// Get full driver details
	return r.GetByID(id)
}

func (r *DriverRepository) Delete(id string) error {
	ctx := context.Background()

	err := r.queries.DeleteDriver(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete driver: %w", err)
	}

	return nil
}

func (r *DriverRepository) CreateOrUpdateLicense(license *domain.DriverLicense) error {
	ctx := context.Background()

	var issuedAt pgtype.Date
	if license.IssuedAt != nil {
		issuedAt = pgtype.Date{Time: *license.IssuedAt, Valid: true}
	}

	var expiresAt pgtype.Date
	if license.ExpiresAt != nil {
		expiresAt = pgtype.Date{Time: *license.ExpiresAt, Valid: true}
	}

	_, err := r.queries.CreateDriverLicense(ctx, sqlc.CreateDriverLicenseParams{
		DriverID:  license.DriverID,
		Number:    license.Number,
		Class:     sqlc.LicenseClass(license.Class),
		IssuedAt:  issuedAt,
		ExpiresAt: expiresAt,
		FileUrl:   license.FileURL,
	})
	if err != nil {
		return fmt.Errorf("failed to create/update driver license: %w", err)
	}

	return nil
}

func (r *DriverRepository) CreateOrUpdateBackgroundCheck(check *domain.DriverBackgroundCheck) error {
	ctx := context.Background()

	var checkedAt pgtype.Timestamptz
	if check.CheckedAt != nil {
		checkedAt = pgtype.Timestamptz{Time: *check.CheckedAt, Valid: true}
	}

	_, err := r.queries.CreateDriverBackgroundCheck(ctx, sqlc.CreateDriverBackgroundCheckParams{
		DriverID:  check.DriverID,
		Status:    sqlc.BackgroundCheckStatus(check.Status),
		FileUrl:   check.FileURL,
		CheckedAt: checkedAt,
	})
	if err != nil {
		return fmt.Errorf("failed to create/update driver background check: %w", err)
	}

	return nil
}

func (r *DriverRepository) CreateOrUpdateAvailability(availability *domain.DriverAvailability) error {
	ctx := context.Background()

	regionsJSON, err := json.Marshal(availability.Regions)
	if err != nil {
		return fmt.Errorf("failed to marshal regions: %w", err)
	}

	daysJSON, err := json.Marshal(availability.Days)
	if err != nil {
		return fmt.Errorf("failed to marshal days: %w", err)
	}

	timeRangesJSON, err := json.Marshal(availability.TimeRanges)
	if err != nil {
		return fmt.Errorf("failed to marshal time ranges: %w", err)
	}

	_, err = r.queries.CreateDriverAvailability(ctx, sqlc.CreateDriverAvailabilityParams{
		DriverID:   availability.DriverID,
		Regions:    regionsJSON,
		Days:       daysJSON,
		TimeRanges: timeRangesJSON,
	})
	if err != nil {
		return fmt.Errorf("failed to create/update driver availability: %w", err)
	}

	return nil
}

func (r *DriverRepository) GetKPIs(driverID string) (*domain.DriverKPIs, error) {
	// For now, return mock data. In a real implementation, this would calculate from reservations and feedback tables
	return &domain.DriverKPIs{
		TotalTrips:    25,
		TotalKM:       1250.5,
		CancelRate:    2.5,
		OnTimeRate:    96.8,
		AverageRating: 4.7,
	}, nil
}

func (r *DriverRepository) mapToDomainDriver(row sqlc.GetDriverByIDRow) *domain.Driver {
	driver := &domain.Driver{
		ID:        row.ID,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		RutOrDNI:  row.RutOrDni,
		Phone:     row.Phone,
		Email:     row.Email,
		PhotoURL:  row.PhotoUrl,
		Status:    domain.DriverStatus(row.Status),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	if row.BirthDate.Valid {
		driver.BirthDate = &row.BirthDate.Time
	}

	// Map license if present
	if row.LicenseNumber != nil {
		license := &domain.DriverLicense{
			DriverID: row.ID,
			Number:   *row.LicenseNumber,
			FileURL:  row.LicenseFileUrl,
		}
		if row.LicenseClass.Valid {
			license.Class = domain.LicenseClass(row.LicenseClass.LicenseClass)
		}
		if row.LicenseIssuedAt.Valid {
			license.IssuedAt = &row.LicenseIssuedAt.Time
		}
		if row.LicenseExpiresAt.Valid {
			license.ExpiresAt = &row.LicenseExpiresAt.Time
		}
		driver.License = license
	}

	// Map background check if present
	if row.BackgroundStatus.Valid {
		check := &domain.DriverBackgroundCheck{
			DriverID: row.ID,
			Status:   domain.BackgroundCheckStatus(row.BackgroundStatus.BackgroundCheckStatus),
			FileURL:  row.BackgroundFileUrl,
		}
		if row.BackgroundCheckedAt.Valid {
			check.CheckedAt = &row.BackgroundCheckedAt.Time
		}
		driver.BackgroundCheck = check
	}

	// Map availability if present
	if row.Regions != nil {
		availability := &domain.DriverAvailability{
			DriverID: row.ID,
		}

		if err := json.Unmarshal(row.Regions, &availability.Regions); err == nil {
			if err := json.Unmarshal(row.Days, &availability.Days); err == nil {
				if err := json.Unmarshal(row.TimeRanges, &availability.TimeRanges); err == nil {
					driver.Availability = availability
				}
			}
		}
	}

	return driver
}

func (r *DriverRepository) mapToDomainDriverFromList(row sqlc.ListDriversRow) *domain.Driver {
	driver := &domain.Driver{
		ID:        row.ID,
		FirstName: row.FirstName,
		LastName:  row.LastName,
		RutOrDNI:  row.RutOrDni,
		Phone:     row.Phone,
		Email:     row.Email,
		PhotoURL:  row.PhotoUrl,
		Status:    domain.DriverStatus(row.Status),
		CreatedAt: row.CreatedAt.Time,
		UpdatedAt: row.UpdatedAt.Time,
	}

	if row.BirthDate.Valid {
		driver.BirthDate = &row.BirthDate.Time
	}

	// Map license if present
	if row.LicenseNumber != nil {
		license := &domain.DriverLicense{
			DriverID: row.ID,
			Number:   *row.LicenseNumber,
		}
		if row.LicenseClass.Valid {
			license.Class = domain.LicenseClass(row.LicenseClass.LicenseClass)
		}
		driver.License = license
	}

	// Map background check if present
	if row.BackgroundStatus.Valid {
		driver.BackgroundCheck = &domain.DriverBackgroundCheck{
			DriverID: row.ID,
			Status:   domain.BackgroundCheckStatus(row.BackgroundStatus.BackgroundCheckStatus),
		}
	}

	return driver
}
