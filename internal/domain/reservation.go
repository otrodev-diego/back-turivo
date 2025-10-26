package domain

import (
	"time"

	"github.com/google/uuid"
)

type ReservationStatus string

const (
	ReservationStatusActiva     ReservationStatus = "ACTIVA"
	ReservationStatusProgramada ReservationStatus = "PROGRAMADA"
	ReservationStatusCompletada ReservationStatus = "COMPLETADA"
	ReservationStatusCancelada  ReservationStatus = "CANCELADA"
)

type Reservation struct {
	ID               string            `json:"id"` // e.g., RSV-1042
	UserID           *uuid.UUID        `json:"user_id,omitempty"`
	OrgID            *uuid.UUID        `json:"org_id,omitempty"`
	Pickup           string            `json:"pickup"`
	Destination      string            `json:"destination"`
	DateTime         time.Time         `json:"datetime"`
	Passengers       int               `json:"passengers"`
	Status           ReservationStatus `json:"status"`
	Amount           *float64          `json:"amount,omitempty"`
	DistanceKM       *float64          `json:"distance_km,omitempty"`
	Notes            *string           `json:"notes,omitempty"`
	AssignedDriverID *string           `json:"assigned_driver_id,omitempty"`
	CreatedAt        time.Time         `json:"created_at"`
	UpdatedAt        time.Time         `json:"updated_at"`

	// Related data
	User           *User            `json:"user,omitempty"`
	AssignedDriver *Driver          `json:"assigned_driver,omitempty"`
	Timeline       []TimelineEvent  `json:"timeline,omitempty"`
	Payments       []Payment        `json:"payments,omitempty"`
	Feedback       []DriverFeedback `json:"feedback,omitempty"`
}

type TimelineEvent struct {
	ID            uuid.UUID `json:"id"`
	ReservationID string    `json:"reservation_id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	At            time.Time `json:"at"`
	Variant       string    `json:"variant"`
	CreatedAt     time.Time `json:"created_at"`
}

type CreateReservationRequest struct {
	UserID      *uuid.UUID `json:"user_id,omitempty"`
	OrgID       *uuid.UUID `json:"org_id,omitempty"`
	Pickup      string     `json:"pickup" validate:"required,min=5,max=500"`
	Destination string     `json:"destination" validate:"required,min=5,max=500"`
	DateTime    time.Time  `json:"datetime" validate:"required"`
	Passengers  int        `json:"passengers" validate:"required,min=1"`
	Notes       *string    `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

type UpdateReservationRequest struct {
	Pickup      *string    `json:"pickup,omitempty" validate:"omitempty,min=5,max=500"`
	Destination *string    `json:"destination,omitempty" validate:"omitempty,min=5,max=500"`
	DateTime    *time.Time `json:"datetime,omitempty"`
	Passengers  *int       `json:"passengers,omitempty" validate:"omitempty,min=1"`
	Amount      *float64   `json:"amount,omitempty"`
	Notes       *string    `json:"notes,omitempty" validate:"omitempty,max=1000"`
}

type ChangeReservationStatusRequest struct {
	NewStatus ReservationStatus `json:"new_status" validate:"required"`
	Notes     *string           `json:"notes,omitempty"`
}

type ListReservationsRequest struct {
	Query    *string            `json:"query,omitempty"`
	Status   *ReservationStatus `json:"status,omitempty"`
	From     *time.Time         `json:"from,omitempty"`
	To       *time.Time         `json:"to,omitempty"`
	UserID   *uuid.UUID         `json:"user_id,omitempty"` // Filter by specific user
	Page     int                `json:"page" validate:"min=1"`
	PageSize int                `json:"page_size" validate:"min=1,max=100"`
	Sort     string             `json:"sort" validate:"omitempty,oneof=datetime created_at"`
}

type CreateTimelineEventRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"required,min=1,max=1000"`
	Variant     string `json:"variant" validate:"omitempty,oneof=default success warning error info"`
}

// Validation methods
func (r *Reservation) CanTransitionTo(newStatus ReservationStatus) bool {
	switch r.Status {
	case ReservationStatusActiva:
		return newStatus == ReservationStatusProgramada || newStatus == ReservationStatusCancelada
	case ReservationStatusProgramada:
		return newStatus == ReservationStatusCompletada || newStatus == ReservationStatusCancelada
	case ReservationStatusCompletada, ReservationStatusCancelada:
		return false // Terminal states
	}
	return false
}

// Price calculation helper
func (r *Reservation) CalculatePrice(vehicleType VehicleType, hasSpecialLanguage bool, stops int) float64 {
	// Base price by vehicle type (aligned with frontend)
	var basePrice float64
	switch vehicleType {
	case VehicleTypeBus:
		basePrice = 150000
	case VehicleTypeVan:
		basePrice = 120000
	case VehicleTypeSedan:
		basePrice = 80000
	case VehicleTypeSUV:
		basePrice = 100000
	default:
		basePrice = 80000
	}

	// Additional passenger cost (>1 passenger)
	if r.Passengers > 1 {
		basePrice += float64(r.Passengers-1) * 5000
	}

	// Additional stops cost
	if stops > 0 {
		basePrice += float64(stops) * 15000
	}

	// Special language cost
	if hasSpecialLanguage {
		basePrice += 25000
	}

	return basePrice
}

// CalculateDistance calculates the distance between pickup and destination
// This is a simplified calculation - in production, you'd use Google Maps API
func (r *Reservation) CalculateDistance() float64 {
	// For now, return a mock distance based on common routes
	// In production, this should call Google Maps Distance Matrix API

	// Mock distance calculation based on pickup and destination
	// This is a simplified version - real implementation would use geocoding
	pickup := r.Pickup
	destination := r.Destination

	// Simple heuristic: if both locations are in the same city, shorter distance
	// This is just for demonstration - real implementation needs proper geocoding
	if len(pickup) > 0 && len(destination) > 0 {
		// Mock calculation: assume average distance of 15-25 km for city trips
		// In production, use Google Maps API to get real distance
		return 18.5 // Mock distance in kilometers
	}

	return 0
}

type ReservationRepository interface {
	Create(reservation *Reservation) error
	GetByID(id string) (*Reservation, error)
	List(req ListReservationsRequest) ([]*Reservation, int, error)
	Update(id string, req UpdateReservationRequest) (*Reservation, error)
	Delete(id string) error
	AssignDriver(id string, driverID string) error
	ChangeStatus(id string, newStatus ReservationStatus) error
	GetTimeline(id string) ([]TimelineEvent, error)
	AddTimelineEvent(id string, event TimelineEvent) error
	GenerateID() string
}
