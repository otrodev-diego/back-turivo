package usecase

import (
	"time"

	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type DriverUseCase struct {
	driverRepo domain.DriverRepository
	logger     *zap.Logger
}

func NewDriverUseCase(driverRepo domain.DriverRepository, logger *zap.Logger) *DriverUseCase {
	return &DriverUseCase{
		driverRepo: driverRepo,
		logger:     logger,
	}
}

func (uc *DriverUseCase) CreateDriver(req domain.CreateDriverRequest) (*domain.Driver, error) {
	uc.logger.Info("Creating driver", zap.String("driver_id", req.ID))

	// Check if driver already exists
	existingDriver, err := uc.driverRepo.GetByID(req.ID)
	if err != nil && err != domain.ErrDriverNotFound {
		uc.logger.Error("Failed to check existing driver", zap.Error(err))
		return nil, domain.ErrInternalError
	}
	if existingDriver != nil {
		return nil, domain.ErrDriverAlreadyExists
	}

	// Create driver
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
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := uc.driverRepo.Create(driver); err != nil {
		uc.logger.Error("Failed to create driver", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	uc.logger.Info("Driver created successfully", zap.String("driver_id", driver.ID))
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

