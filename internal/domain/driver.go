package domain

import (
	"encoding/json"
	"time"
)

type DriverStatus string

const (
	DriverStatusActive   DriverStatus = "ACTIVE"
	DriverStatusInactive DriverStatus = "INACTIVE"
)

type LicenseClass string

const (
	LicenseClassA1 LicenseClass = "A1"
	LicenseClassA2 LicenseClass = "A2"
	LicenseClassA3 LicenseClass = "A3"
	LicenseClassA4 LicenseClass = "A4"
	LicenseClassA5 LicenseClass = "A5"
	LicenseClassB  LicenseClass = "B"
	LicenseClassC  LicenseClass = "C"
	LicenseClassD  LicenseClass = "D"
	LicenseClassE  LicenseClass = "E"
)

type BackgroundCheckStatus string

const (
	BackgroundCheckStatusApproved BackgroundCheckStatus = "APPROVED"
	BackgroundCheckStatusPending  BackgroundCheckStatus = "PENDING"
	BackgroundCheckStatusRejected BackgroundCheckStatus = "REJECTED"
)

type VehicleType string

const (
	VehicleTypeBus   VehicleType = "BUS"
	VehicleTypeVan   VehicleType = "VAN"
	VehicleTypeSedan VehicleType = "SEDAN"
	VehicleTypeSUV   VehicleType = "SUV"
)

type TimeRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type Driver struct {
	ID        string       `json:"id"`
	FirstName string       `json:"first_name"`
	LastName  string       `json:"last_name"`
	RutOrDNI  string       `json:"rut_or_dni"`
	BirthDate *time.Time   `json:"birth_date,omitempty"`
	Phone     *string      `json:"phone,omitempty"`
	Email     *string      `json:"email,omitempty"`
	PhotoURL  *string      `json:"photo_url,omitempty"`
	Status    DriverStatus `json:"status"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`

	// New relationships
	UserID    *string `json:"user_id,omitempty"`
	CompanyID *string `json:"company_id,omitempty"`
	VehicleID *string `json:"vehicle_id,omitempty"`

	// Related data
	License         *DriverLicense         `json:"license,omitempty"`
	BackgroundCheck *DriverBackgroundCheck `json:"background_check,omitempty"`
	Availability    *DriverAvailability    `json:"availability,omitempty"`
	Vehicle         *Vehicle               `json:"vehicle,omitempty"`
	KPIs            *DriverKPIs            `json:"kpis,omitempty"`
	User            *User                  `json:"user,omitempty"`
	Company         *Company               `json:"company,omitempty"`
}

type DriverLicense struct {
	DriverID  string       `json:"driver_id"`
	Number    string       `json:"number"`
	Class     LicenseClass `json:"class"`
	IssuedAt  *time.Time   `json:"issued_at,omitempty"`
	ExpiresAt *time.Time   `json:"expires_at,omitempty"`
	FileURL   *string      `json:"file_url,omitempty"`
}

type DriverBackgroundCheck struct {
	DriverID  string                `json:"driver_id"`
	Status    BackgroundCheckStatus `json:"status"`
	FileURL   *string               `json:"file_url,omitempty"`
	CheckedAt *time.Time            `json:"checked_at,omitempty"`
}

type DriverAvailability struct {
	DriverID   string      `json:"driver_id"`
	Regions    []string    `json:"regions"`
	Days       []string    `json:"days"`
	TimeRanges []TimeRange `json:"time_ranges"`
	UpdatedAt  time.Time   `json:"updated_at"`
}

type DriverKPIs struct {
	TotalTrips    int     `json:"total_trips"`
	TotalKM       float64 `json:"total_km"`
	CancelRate    float64 `json:"cancel_rate"`
	OnTimeRate    float64 `json:"on_time_rate"`
	AverageRating float64 `json:"average_rating"`
}

type DriverFeedback struct {
	ID            string    `json:"id"`
	DriverID      string    `json:"driver_id"`
	ReservationID string    `json:"reservation_id"`
	Rating        float64   `json:"rating"`
	Comment       *string   `json:"comment,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateDriverRequest struct {
	ID        string       `json:"id" validate:"required,min=3,max=20"`
	FirstName string       `json:"first_name" validate:"required,min=2,max=255"`
	LastName  string       `json:"last_name" validate:"required,min=2,max=255"`
	RutOrDNI  string       `json:"rut_or_dni" validate:"required,min=8,max=50"`
	BirthDate *time.Time   `json:"birth_date,omitempty"`
	Phone     *string      `json:"phone,omitempty" validate:"omitempty,min=8,max=50"`
	Email     *string      `json:"email,omitempty" validate:"omitempty,email"`
	PhotoURL  *string      `json:"photo_url,omitempty" validate:"omitempty,url"`
	Status    DriverStatus `json:"status" validate:"required"`

	// New fields for relationships
	UserID    *string `json:"user_id,omitempty"`
	CompanyID *string `json:"company_id,omitempty"`
	VehicleID *string `json:"vehicle_id,omitempty"`
}

type UpdateDriverRequest struct {
	FirstName *string       `json:"first_name,omitempty" validate:"omitempty,min=2,max=255"`
	LastName  *string       `json:"last_name,omitempty" validate:"omitempty,min=2,max=255"`
	RutOrDNI  *string       `json:"rut_or_dni,omitempty" validate:"omitempty,min=8,max=50"`
	BirthDate *time.Time    `json:"birth_date,omitempty"`
	Phone     *string       `json:"phone,omitempty" validate:"omitempty,min=8,max=50"`
	Email     *string       `json:"email,omitempty" validate:"omitempty,email"`
	PhotoURL  *string       `json:"photo_url,omitempty" validate:"omitempty,url"`
	Status    *DriverStatus `json:"status,omitempty"`

	// New fields for relationships
	UserID    *string `json:"user_id,omitempty"`
	CompanyID *string `json:"company_id,omitempty"`
	VehicleID *string `json:"vehicle_id,omitempty"`
}

type ListDriversRequest struct {
	Query    *string       `json:"query,omitempty"`
	Status   *DriverStatus `json:"status,omitempty"`
	Page     int           `json:"page" validate:"min=1"`
	PageSize int           `json:"page_size" validate:"min=1,max=100"`
	Sort     string        `json:"sort" validate:"omitempty,oneof=name id created_at"`
}

// Helper methods for JSON marshaling/unmarshaling of availability data
func (da *DriverAvailability) MarshalRegions() ([]byte, error) {
	return json.Marshal(da.Regions)
}

func (da *DriverAvailability) UnmarshalRegions(data []byte) error {
	return json.Unmarshal(data, &da.Regions)
}

func (da *DriverAvailability) MarshalDays() ([]byte, error) {
	return json.Marshal(da.Days)
}

func (da *DriverAvailability) UnmarshalDays(data []byte) error {
	return json.Unmarshal(data, &da.Days)
}

func (da *DriverAvailability) MarshalTimeRanges() ([]byte, error) {
	return json.Marshal(da.TimeRanges)
}

func (da *DriverAvailability) UnmarshalTimeRanges(data []byte) error {
	return json.Unmarshal(data, &da.TimeRanges)
}

type DriverRepository interface {
	Create(driver *Driver) error
	GetByID(id string) (*Driver, error)
	List(req ListDriversRequest) ([]*Driver, int, error)
	Update(id string, req UpdateDriverRequest) (*Driver, error)
	Delete(id string) error

	// License operations
	CreateOrUpdateLicense(license *DriverLicense) error

	// Background check operations
	CreateOrUpdateBackgroundCheck(check *DriverBackgroundCheck) error

	// Availability operations
	CreateOrUpdateAvailability(availability *DriverAvailability) error

	// KPIs (calculated from other tables)
	GetKPIs(driverID string) (*DriverKPIs, error)

	// Feedback operations
	CreateFeedback(feedback *DriverFeedback) error
	GetDriverFeedback(driverID string) ([]*DriverFeedback, error)

	// New methods for driver dashboard
	GetByUserID(userID string) (*Driver, error)
	GetDriverTrips(driverID string) ([]*Reservation, error)
	GetDriverVehicle(driverID string) (*Vehicle, error)
	UpdateTripStatus(driverID, tripID, status string) error
}
