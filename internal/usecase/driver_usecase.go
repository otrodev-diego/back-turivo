package usecase

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type DriverUseCase struct {
	driverRepo            domain.DriverRepository
	userRepo              domain.UserRepository
	emailService          domain.EmailService
	registrationTokenRepo domain.RegistrationTokenRepository
	passwordService       domain.PasswordService
	logger                *zap.Logger
}

func NewDriverUseCase(driverRepo domain.DriverRepository, userRepo domain.UserRepository, emailService domain.EmailService, registrationTokenRepo domain.RegistrationTokenRepository, passwordService domain.PasswordService, logger *zap.Logger) *DriverUseCase {
	return &DriverUseCase{
		driverRepo:            driverRepo,
		userRepo:              userRepo,
		emailService:          emailService,
		registrationTokenRepo: registrationTokenRepo,
		passwordService:       passwordService,
		logger:                logger,
	}
}

func (uc *DriverUseCase) CreateDriver(req domain.CreateDriverRequest) (*domain.Driver, error) {
	uc.logger.Info("🚀 === CreateDriver UseCase Started ===",
		zap.String("driver_id", req.ID),
		zap.String("first_name", req.FirstName),
		zap.String("last_name", req.LastName),
		zap.String("email", *req.Email),
		zap.String("rut_or_dni", req.RutOrDNI),
	)

	// Check if driver already exists
	uc.logger.Info("🔍 Checking if driver already exists", zap.String("driver_id", req.ID))
	existingDriver, err := uc.driverRepo.GetByID(req.ID)
	if err != nil && err != domain.ErrDriverNotFound {
		uc.logger.Error("❌ FAILED to check existing driver", zap.Error(err))
		return nil, domain.ErrInternalError
	}
	if existingDriver != nil {
		uc.logger.Warn("⚠️ Driver already exists", zap.String("driver_id", req.ID))
		return nil, domain.ErrDriverAlreadyExists
	}
	uc.logger.Info("✅ Driver doesn't exist, proceeding", zap.String("driver_id", req.ID))

	// Create driver
	uc.logger.Info("📝 Creating driver object", zap.String("driver_id", req.ID))
	driver := &domain.Driver{
		ID:        req.ID,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		RutOrDNI:  req.RutOrDNI,
		BirthDate: req.BirthDate,
		Phone:     req.Phone,
		Email:     req.Email,
		PhotoURL:  req.PhotoURL,
		Status:    req.Status,
		UserID:    req.UserID,
		CompanyID: req.CompanyID,
		VehicleID: req.VehicleID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	uc.logger.Info("✅ Driver object created", zap.String("driver_id", driver.ID))

	// If email is provided, create a user account for the driver
	if req.Email != nil && *req.Email != "" {
		uc.logger.Info("📧 === Creating user account for driver ===", zap.String("email", *req.Email))

		// Check if user already exists
		uc.logger.Info("🔍 Checking if user already exists", zap.String("email", *req.Email))
		existingUser, err := uc.userRepo.GetByEmail(*req.Email)
		if err != nil {
			uc.logger.Info("🔍 User lookup result",
				zap.String("email", *req.Email),
				zap.Error(err),
				zap.String("error_type", fmt.Sprintf("%T", err)),
				zap.Bool("is_user_not_found", err == domain.ErrUserNotFound),
				zap.String("error_string", err.Error()),
			)
			// If it's not "user not found", it's a real error
			if err != domain.ErrUserNotFound {
				uc.logger.Error("❌ FAILED to check existing user", zap.Error(err))
				return nil, domain.ErrInternalError
			}
			// If it's "user not found", that's expected - user doesn't exist
			uc.logger.Info("✅ User not found (expected) - will create new user", zap.String("email", *req.Email))
		}

		if existingUser == nil {
			uc.logger.Info("✅ User doesn't exist, creating new user", zap.String("email", *req.Email))

			// Generate a temporary password for the driver
			tempPassword := "TempPassword123!" // This will be changed when driver completes registration

			// Hash the temporary password
			hashedPassword, err := uc.passwordService.HashPassword(tempPassword)
			if err != nil {
				uc.logger.Error("❌ FAILED to hash temporary password", zap.Error(err))
				return nil, domain.ErrInternalError
			}

			// Create user with DRIVER role
			user := &domain.User{
				ID:           uuid.New(), // Generate UUID before creating
				Name:         req.FirstName + " " + req.LastName,
				Email:        *req.Email,
				PasswordHash: hashedPassword, // Use real hashed password
				Role:         domain.UserRoleDriver,
				Status:       domain.UserStatusActive,
			}

			uc.logger.Info("📝 Creating user in database", zap.String("user_id", user.ID.String()))
			err = uc.userRepo.Create(user)
			if err != nil {
				uc.logger.Error("❌ FAILED to create user for driver", zap.Error(err))
				return nil, domain.ErrInternalError
			}

			// Link user to driver
			userIDStr := user.ID.String()
			driver.UserID = &userIDStr
			uc.logger.Info("✅ User created and linked to driver", zap.String("user_id", user.ID.String()))
		} else {
			uc.logger.Info("⚠️ User already exists, linking to driver", zap.String("email", *req.Email))
			// User already exists, link to driver
			userIDStr := existingUser.ID.String()
			driver.UserID = &userIDStr
			uc.logger.Info("✅ Existing user linked to driver", zap.String("user_id", existingUser.ID.String()))
		}
	}

	uc.logger.Info("💾 === Creating driver in database ===", zap.String("driver_id", driver.ID))
	if err := uc.driverRepo.Create(driver); err != nil {
		uc.logger.Error("❌ FAILED to create driver in database", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	uc.logger.Info("🎉 === Driver created successfully ===", zap.String("driver_id", driver.ID))

	// Send welcome email to driver
	if req.Email != nil && *req.Email != "" {
		uc.logger.Info("📧 === Sending welcome email to driver ===", zap.String("email", *req.Email))
		name := req.FirstName + " " + req.LastName

		// Generate a proper registration token
		registrationToken := &domain.RegistrationToken{
			ID:        uuid.New(),
			Token:     uuid.New().String(), // Generate a proper token
			Email:     *req.Email,
			OrgID:     nil, // Driver doesn't belong to an organization initially
			Role:      domain.UserRoleDriver,
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // 24 hours expiration
			Used:      false,
			CreatedAt: time.Now().Unix(),
		}

		// Save the registration token to database
		uc.logger.Info("💾 Saving registration token to database", zap.String("token", registrationToken.Token))
		err := uc.registrationTokenRepo.Create(registrationToken)
		if err != nil {
			uc.logger.Error("❌ FAILED to save registration token", zap.Error(err))
			// Don't return error, just log it - the driver was created successfully
		} else {
			uc.logger.Info("✅ Registration token saved successfully", zap.String("token", registrationToken.Token))
		}

		token := registrationToken.Token

		emailErr := uc.emailService.SendWelcomeEmail(*req.Email, name, token)
		if emailErr != nil {
			uc.logger.Error("❌ FAILED to send welcome email", zap.Error(emailErr))
			// Don't return error, just log it - the driver was created successfully
		} else {
			uc.logger.Info("✅ Welcome email sent successfully", zap.String("email", *req.Email))
		}
	}

	return driver, nil
}

func (uc *DriverUseCase) GetDriverByID(id string) (*domain.Driver, error) {
	driver, err := uc.driverRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrDriverNotFound
		}
		uc.logger.Error("Failed to get driver by ID", zap.Error(err), zap.String("driver_id", id))
		return nil, domain.ErrInternalError
	}

	return driver, nil
}

func (uc *DriverUseCase) ListDrivers(req domain.ListDriversRequest) ([]*domain.Driver, int, error) {
	// Set default pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}
	if req.Sort == "" {
		req.Sort = "created_at"
	}

	drivers, total, err := uc.driverRepo.List(req)
	if err != nil {
		uc.logger.Error("Failed to list drivers", zap.Error(err))
		return nil, 0, domain.ErrInternalError
	}

	return drivers, total, nil
}

func (uc *DriverUseCase) UpdateDriver(id string, req domain.UpdateDriverRequest) (*domain.Driver, error) {
	uc.logger.Info("Updating driver", zap.String("driver_id", id))

	// Check if driver exists
	_, err := uc.driverRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrDriverNotFound
		}
		uc.logger.Error("Failed to get driver for update", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	driver, err := uc.driverRepo.Update(id, req)
	if err != nil {
		uc.logger.Error("Failed to update driver", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	uc.logger.Info("Driver updated successfully", zap.String("driver_id", driver.ID))
	return driver, nil
}

func (uc *DriverUseCase) GetDriverByUserID(userID string) (*domain.Driver, error) {
	uc.logger.Info("🔍 Getting driver by user ID", zap.String("user_id", userID))

	driver, err := uc.driverRepo.GetByUserID(userID)
	if err != nil {
		uc.logger.Error("❌ Failed to get driver by user ID", zap.Error(err), zap.String("user_id", userID))
		return nil, err
	}

	uc.logger.Info("✅ Driver found by user ID", zap.String("driver_id", driver.ID), zap.String("user_id", userID))
	return driver, nil
}

func (uc *DriverUseCase) DeleteDriver(id string) error {
	uc.logger.Info("Deleting driver", zap.String("driver_id", id))

	// Check if driver exists
	_, err := uc.driverRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrDriverNotFound
		}
		uc.logger.Error("Failed to get driver for deletion", zap.Error(err))
		return domain.ErrInternalError
	}

	if err := uc.driverRepo.Delete(id); err != nil {
		uc.logger.Error("Failed to delete driver", zap.Error(err))
		return domain.ErrInternalError
	}

	uc.logger.Info("Driver deleted successfully", zap.String("driver_id", id))
	return nil
}

func (uc *DriverUseCase) UpdateDriverLicense(driverID string, license domain.DriverLicense) error {
	uc.logger.Info("Updating driver license", zap.String("driver_id", driverID))

	// Check if driver exists
	_, err := uc.driverRepo.GetByID(driverID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrDriverNotFound
		}
		uc.logger.Error("Failed to get driver for license update", zap.Error(err))
		return domain.ErrInternalError
	}

	license.DriverID = driverID
	if err := uc.driverRepo.CreateOrUpdateLicense(&license); err != nil {
		uc.logger.Error("Failed to update driver license", zap.Error(err))
		return domain.ErrInternalError
	}

	uc.logger.Info("Driver license updated successfully", zap.String("driver_id", driverID))
	return nil
}

func (uc *DriverUseCase) UpdateDriverBackgroundCheck(driverID string, check domain.DriverBackgroundCheck) error {
	uc.logger.Info("Updating driver background check", zap.String("driver_id", driverID))

	// Check if driver exists
	_, err := uc.driverRepo.GetByID(driverID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrDriverNotFound
		}
		uc.logger.Error("Failed to get driver for background check update", zap.Error(err))
		return domain.ErrInternalError
	}

	check.DriverID = driverID
	if check.Status == domain.BackgroundCheckStatusApproved && check.CheckedAt == nil {
		now := time.Now()
		check.CheckedAt = &now
	}

	if err := uc.driverRepo.CreateOrUpdateBackgroundCheck(&check); err != nil {
		uc.logger.Error("Failed to update driver background check", zap.Error(err))
		return domain.ErrInternalError
	}

	uc.logger.Info("Driver background check updated successfully", zap.String("driver_id", driverID))
	return nil
}

func (uc *DriverUseCase) UpdateDriverAvailability(driverID string, availability domain.DriverAvailability) error {
	uc.logger.Info("Updating driver availability", zap.String("driver_id", driverID))

	// Check if driver exists
	_, err := uc.driverRepo.GetByID(driverID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrDriverNotFound
		}
		uc.logger.Error("Failed to get driver for availability update", zap.Error(err))
		return domain.ErrInternalError
	}

	availability.DriverID = driverID
	availability.UpdatedAt = time.Now()

	if err := uc.driverRepo.CreateOrUpdateAvailability(&availability); err != nil {
		uc.logger.Error("Failed to update driver availability", zap.Error(err))
		return domain.ErrInternalError
	}

	uc.logger.Info("Driver availability updated successfully", zap.String("driver_id", driverID))
	return nil
}

func (uc *DriverUseCase) GetDriverKPIs(driverID string) (*domain.DriverKPIs, error) {
	// Check if driver exists
	_, err := uc.driverRepo.GetByID(driverID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrDriverNotFound
		}
		uc.logger.Error("Failed to get driver for KPIs", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	kpis, err := uc.driverRepo.GetKPIs(driverID)
	if err != nil {
		uc.logger.Error("Failed to get driver KPIs", zap.Error(err), zap.String("driver_id", driverID))
		return nil, domain.ErrInternalError
	}

	return kpis, nil
}

// GetDriverTrips gets trips for a specific driver
func (uc *DriverUseCase) GetDriverTrips(driverID string) ([]*domain.Reservation, error) {
	uc.logger.Info("Getting driver trips", zap.String("driver_id", driverID))

	trips, err := uc.driverRepo.GetDriverTrips(driverID)
	if err != nil {
		uc.logger.Error("Failed to get driver trips", zap.Error(err), zap.String("driver_id", driverID))
		return nil, domain.ErrInternalError
	}

	// If no trips found, return empty slice
	if trips == nil {
		trips = []*domain.Reservation{}
	}

	uc.logger.Info("Driver trips retrieved", zap.String("driver_id", driverID), zap.Int("count", len(trips)))
	return trips, nil
}

// GetDriverVehicle gets the vehicle assigned to a driver
func (uc *DriverUseCase) GetDriverVehicle(driverID string) (*domain.Vehicle, error) {
	uc.logger.Info("Getting driver vehicle", zap.String("driver_id", driverID))

	vehicle, err := uc.driverRepo.GetDriverVehicle(driverID)
	if err != nil {
		uc.logger.Error("Failed to get driver vehicle", zap.Error(err))
		return nil, err
	}

	uc.logger.Info("Driver vehicle retrieved", zap.String("driver_id", driverID))
	return vehicle, nil
}

// UpdateTripStatus updates the status of a trip
func (uc *DriverUseCase) UpdateTripStatus(driverID, tripID, status string) error {
	uc.logger.Info("Updating trip status", zap.String("driver_id", driverID), zap.String("trip_id", tripID), zap.String("status", status))

	err := uc.driverRepo.UpdateTripStatus(driverID, tripID, status)
	if err != nil {
		uc.logger.Error("Failed to update trip status", zap.Error(err))
		return domain.ErrInternalError
	}

	uc.logger.Info("Trip status updated successfully", zap.String("driver_id", driverID), zap.String("trip_id", tripID))
	return nil
}
