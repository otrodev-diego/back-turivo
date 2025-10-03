package domain

import (
	"time"
)

type VehicleStatus string

const (
	VehicleStatusAvailable   VehicleStatus = "AVAILABLE"
	VehicleStatusAssigned    VehicleStatus = "ASSIGNED"
	VehicleStatusMaintenance VehicleStatus = "MAINTENANCE"
	VehicleStatusInactive    VehicleStatus = "INACTIVE"
)

// Vehicle represents a vehicle in the system
type Vehicle struct {
	ID                  string        `json:"id"`
	DriverID            *string       `json:"driver_id,omitempty"`
	Type                VehicleType   `json:"type"`
	Brand               string        `json:"brand"`
	Model               string        `json:"model"`
	Year                *int          `json:"year,omitempty"`
	Plate               *string       `json:"plate,omitempty"`
	VIN                 *string       `json:"vin,omitempty"`
	Color               *string       `json:"color,omitempty"`
	Capacity            *int          `json:"capacity,omitempty"` // Passenger capacity
	InsurancePolicy     *string       `json:"insurance_policy,omitempty"`
	InsuranceExpiresAt  *time.Time    `json:"insurance_expires_at,omitempty"`
	InspectionExpiresAt *time.Time    `json:"inspection_expires_at,omitempty"`
	Status              VehicleStatus `json:"status"`
	Photos              []string      `json:"photos,omitempty"`
	CreatedAt           time.Time     `json:"created_at"`
	UpdatedAt           time.Time     `json:"updated_at"`

	// Related data
	Driver *DriverBasicInfo `json:"driver,omitempty"`
}

// DriverBasicInfo contains basic driver information for vehicle responses
type DriverBasicInfo struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Phone     string `json:"phone,omitempty"`
}

type CreateVehicleRequest struct {
	Type                VehicleType   `json:"type" validate:"required"`
	Brand               string        `json:"brand" validate:"required,min=2,max=100"`
	Model               string        `json:"model" validate:"required,min=2,max=100"`
	Year                *int          `json:"year,omitempty" validate:"omitempty,min=1900,max=2100"`
	Plate               *string       `json:"plate,omitempty" validate:"omitempty,min=5,max=20"`
	VIN                 *string       `json:"vin,omitempty" validate:"omitempty,min=10,max=100"`
	Color               *string       `json:"color,omitempty" validate:"omitempty,max=50"`
	Capacity            *int          `json:"capacity,omitempty" validate:"omitempty,min=1,max=60"`
	InsurancePolicy     *string       `json:"insurance_policy,omitempty" validate:"omitempty,max=100"`
	InsuranceExpiresAt  *time.Time    `json:"insurance_expires_at,omitempty"`
	InspectionExpiresAt *time.Time    `json:"inspection_expires_at,omitempty"`
	Status              VehicleStatus `json:"status" validate:"required"`
}

type UpdateVehicleRequest struct {
	Type                *VehicleType   `json:"type,omitempty"`
	Brand               *string        `json:"brand,omitempty" validate:"omitempty,min=2,max=100"`
	Model               *string        `json:"model,omitempty" validate:"omitempty,min=2,max=100"`
	Year                *int           `json:"year,omitempty" validate:"omitempty,min=1900,max=2100"`
	Plate               *string        `json:"plate,omitempty" validate:"omitempty,min=5,max=20"`
	VIN                 *string        `json:"vin,omitempty" validate:"omitempty,min=10,max=100"`
	Color               *string        `json:"color,omitempty" validate:"omitempty,max=50"`
	Capacity            *int           `json:"capacity,omitempty" validate:"omitempty,min=1,max=60"`
	InsurancePolicy     *string        `json:"insurance_policy,omitempty" validate:"omitempty,max=100"`
	InsuranceExpiresAt  *time.Time     `json:"insurance_expires_at,omitempty"`
	InspectionExpiresAt *time.Time     `json:"inspection_expires_at,omitempty"`
	Status              *VehicleStatus `json:"status,omitempty"`
}

type AssignVehicleRequest struct {
	VehicleID string  `json:"vehicle_id" validate:"required,uuid"`
	DriverID  *string `json:"driver_id,omitempty"` // nil to unassign
}

type ListVehiclesRequest struct {
	Query    *string        `json:"query,omitempty"`
	Type     *VehicleType   `json:"type,omitempty"`
	Status   *VehicleStatus `json:"status,omitempty"`
	DriverID *string        `json:"driver_id,omitempty"`
	Page     int            `json:"page" validate:"min=1"`
	PageSize int            `json:"page_size" validate:"min=1,max=100"`
	Sort     string         `json:"sort" validate:"omitempty,oneof=brand model year created_at"`
}

type VehicleRepository interface {
	Create(vehicle *Vehicle) error
	GetByID(id string) (*Vehicle, error)
	GetByDriverID(driverID string) (*Vehicle, error)
	List(req ListVehiclesRequest) ([]*Vehicle, int, error)
	Update(id string, req UpdateVehicleRequest) (*Vehicle, error)
	Delete(id string) error
	AssignToDriver(vehicleID string, driverID *string) error
	
	// Photo operations
	AddPhoto(vehicleID string, photoURL string) error
	RemovePhoto(photoID string) error
	GetPhotos(vehicleID string) ([]string, error)
}

type VehicleUseCase interface {
	CreateVehicle(req CreateVehicleRequest) (*Vehicle, error)
	GetVehicle(id string) (*Vehicle, error)
	ListVehicles(req ListVehiclesRequest) ([]*Vehicle, int, error)
	UpdateVehicle(id string, req UpdateVehicleRequest) (*Vehicle, error)
	DeleteVehicle(id string) error
	AssignVehicleToDriver(vehicleID string, driverID *string) (*Vehicle, error)
}

