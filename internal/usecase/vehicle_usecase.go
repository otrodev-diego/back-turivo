package usecase

import (
	"fmt"
	"time"

	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type vehicleUseCase struct {
	vehicleRepo domain.VehicleRepository
	driverRepo  domain.DriverRepository
	logger      *zap.Logger
}

func NewVehicleUseCase(
	vehicleRepo domain.VehicleRepository,
	driverRepo domain.DriverRepository,
	logger *zap.Logger,
) domain.VehicleUseCase {
	return &vehicleUseCase{
		vehicleRepo: vehicleRepo,
		driverRepo:  driverRepo,
		logger:      logger,
	}
}

func (uc *vehicleUseCase) CreateVehicle(req domain.CreateVehicleRequest) (*domain.Vehicle, error) {
	uc.logger.Info("Creating vehicle",
		zap.String("brand", req.Brand),
		zap.String("model", req.Model),
		zap.String("type", string(req.Type)))

	// Validate vehicle type
	validTypes := map[domain.VehicleType]bool{
		domain.VehicleTypeBus:   true,
		domain.VehicleTypeVan:   true,
		domain.VehicleTypeSedan: true,
		domain.VehicleTypeSUV:   true,
	}
	if !validTypes[req.Type] {
		return nil, domain.ErrInvalidInput
	}

	// Validate status
	validStatuses := map[domain.VehicleStatus]bool{
		domain.VehicleStatusAvailable:   true,
		domain.VehicleStatusAssigned:    true,
		domain.VehicleStatusMaintenance: true,
		domain.VehicleStatusInactive:    true,
	}
	if !validStatuses[req.Status] {
		return nil, domain.ErrInvalidInput
	}

	vehicle := &domain.Vehicle{
		Type:                req.Type,
		Brand:               req.Brand,
		Model:               req.Model,
		Year:                req.Year,
		Plate:               req.Plate,
		VIN:                 req.VIN,
		Color:               req.Color,
		Capacity:            req.Capacity,
		InsurancePolicy:     req.InsurancePolicy,
		InsuranceExpiresAt:  req.InsuranceExpiresAt,
		InspectionExpiresAt: req.InspectionExpiresAt,
		Status:              req.Status,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	if err := uc.vehicleRepo.Create(vehicle); err != nil {
		uc.logger.Error("Failed to create vehicle", zap.Error(err))
		return nil, fmt.Errorf("failed to create vehicle: %w", err)
	}

	uc.logger.Info("Vehicle created successfully", zap.String("id", vehicle.ID))
	return vehicle, nil
}

func (uc *vehicleUseCase) GetVehicle(id string) (*domain.Vehicle, error) {
	uc.logger.Info("Getting vehicle", zap.String("id", id))

	vehicle, err := uc.vehicleRepo.GetByID(id)
	if err != nil {
		uc.logger.Error("Failed to get vehicle", zap.Error(err), zap.String("id", id))
		return nil, err
	}

	return vehicle, nil
}

func (uc *vehicleUseCase) ListVehicles(req domain.ListVehiclesRequest) ([]*domain.Vehicle, int, error) {
	uc.logger.Info("Listing vehicles",
		zap.Int("page", req.Page),
		zap.Int("page_size", req.PageSize))

	// Set defaults
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 10
	}

	vehicles, total, err := uc.vehicleRepo.List(req)
	if err != nil {
		uc.logger.Error("Failed to list vehicles", zap.Error(err))
		return nil, 0, fmt.Errorf("failed to list vehicles: %w", err)
	}

	uc.logger.Info("Vehicles listed successfully",
		zap.Int("total", total),
		zap.Int("returned", len(vehicles)))

	return vehicles, total, nil
}

func (uc *vehicleUseCase) UpdateVehicle(id string, req domain.UpdateVehicleRequest) (*domain.Vehicle, error) {
	uc.logger.Info("Updating vehicle", zap.String("id", id))

	// Check if vehicle exists
	_, err := uc.vehicleRepo.GetByID(id)
	if err != nil {
		uc.logger.Error("Vehicle not found", zap.Error(err), zap.String("id", id))
		return nil, err
	}

	// Validate vehicle type if provided
	if req.Type != nil {
		validTypes := map[domain.VehicleType]bool{
			domain.VehicleTypeBus:   true,
			domain.VehicleTypeVan:   true,
			domain.VehicleTypeSedan: true,
			domain.VehicleTypeSUV:   true,
		}
		if !validTypes[*req.Type] {
			return nil, domain.ErrInvalidInput
		}
	}

	// Validate status if provided
	if req.Status != nil {
		validStatuses := map[domain.VehicleStatus]bool{
			domain.VehicleStatusAvailable:   true,
			domain.VehicleStatusAssigned:    true,
			domain.VehicleStatusMaintenance: true,
			domain.VehicleStatusInactive:    true,
		}
		if !validStatuses[*req.Status] {
			return nil, domain.ErrInvalidInput
		}
	}

	vehicle, err := uc.vehicleRepo.Update(id, req)
	if err != nil {
		uc.logger.Error("Failed to update vehicle", zap.Error(err), zap.String("id", id))
		return nil, fmt.Errorf("failed to update vehicle: %w", err)
	}

	uc.logger.Info("Vehicle updated successfully", zap.String("id", id))
	return vehicle, nil
}

func (uc *vehicleUseCase) DeleteVehicle(id string) error {
	uc.logger.Info("Deleting vehicle", zap.String("id", id))

	// Check if vehicle exists
	vehicle, err := uc.vehicleRepo.GetByID(id)
	if err != nil {
		uc.logger.Error("Vehicle not found", zap.Error(err), zap.String("id", id))
		return err
	}

	// Don't allow deletion if vehicle is assigned to a driver
	if vehicle.DriverID != nil {
		uc.logger.Warn("Cannot delete vehicle assigned to a driver",
			zap.String("vehicle_id", id),
			zap.String("driver_id", *vehicle.DriverID))
		return fmt.Errorf("cannot delete vehicle assigned to a driver: unassign first")
	}

	if err := uc.vehicleRepo.Delete(id); err != nil {
		uc.logger.Error("Failed to delete vehicle", zap.Error(err), zap.String("id", id))
		return fmt.Errorf("failed to delete vehicle: %w", err)
	}

	uc.logger.Info("Vehicle deleted successfully", zap.String("id", id))
	return nil
}

func (uc *vehicleUseCase) AssignVehicleToDriver(vehicleID string, driverID *string) (*domain.Vehicle, error) {
	uc.logger.Info("Assigning vehicle to driver",
		zap.String("vehicle_id", vehicleID),
		zap.Any("driver_id", driverID))

	// Check if vehicle exists
	_, err := uc.vehicleRepo.GetByID(vehicleID)
	if err != nil {
		uc.logger.Error("Vehicle not found", zap.Error(err), zap.String("id", vehicleID))
		return nil, err
	}

	// If assigning to a driver, check if driver exists
	if driverID != nil {
		_, err := uc.driverRepo.GetByID(*driverID)
		if err != nil {
			uc.logger.Error("Driver not found", zap.Error(err), zap.String("id", *driverID))
			return nil, fmt.Errorf("driver not found: %w", err)
		}
	}

	// Perform the assignment
	if err := uc.vehicleRepo.AssignToDriver(vehicleID, driverID); err != nil {
		uc.logger.Error("Failed to assign vehicle", zap.Error(err))
		return nil, fmt.Errorf("failed to assign vehicle: %w", err)
	}

	// Return updated vehicle
	vehicle, err := uc.vehicleRepo.GetByID(vehicleID)
	if err != nil {
		uc.logger.Error("Failed to get updated vehicle", zap.Error(err))
		return nil, err
	}

	if driverID != nil {
		uc.logger.Info("Vehicle assigned successfully",
			zap.String("vehicle_id", vehicleID),
			zap.String("driver_id", *driverID))
	} else {
		uc.logger.Info("Vehicle unassigned successfully",
			zap.String("vehicle_id", vehicleID))
	}

	return vehicle, nil
}

