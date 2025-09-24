package domain

import (
	"time"

	"github.com/google/uuid"
)

type DriverFeedback struct {
	ID            uuid.UUID `json:"id"`
	TripID        string    `json:"trip_id"`
	PassengerName string    `json:"passenger_name"`
	Rating        int       `json:"rating"`
	Comment       *string   `json:"comment,omitempty"`
	CreatedAt     time.Time `json:"created_at"`

	// Related data
	Trip *Reservation `json:"trip,omitempty"`
}

type CreateFeedbackRequest struct {
	TripID        string  `json:"trip_id" validate:"required"`
	PassengerName string  `json:"passenger_name" validate:"required,min=2,max=255"`
	Rating        int     `json:"rating" validate:"required,min=1,max=5"`
	Comment       *string `json:"comment,omitempty" validate:"omitempty,max=1000"`
}

type ListFeedbackRequest struct {
	TripID   *string `json:"trip_id,omitempty"`
	DriverID *string `json:"driver_id,omitempty"`
	Page     int     `json:"page" validate:"min=1"`
	PageSize int     `json:"page_size" validate:"min=1,max=100"`
	Sort     string  `json:"sort" validate:"omitempty,oneof=rating created_at"`
}

type FeedbackRepository interface {
	Create(feedback *DriverFeedback) error
	GetByID(id uuid.UUID) (*DriverFeedback, error)
	GetByTripID(tripID string) ([]*DriverFeedback, error)
	GetByDriverID(driverID string) ([]*DriverFeedback, error)
	List(req ListFeedbackRequest) ([]*DriverFeedback, int, error)
	Delete(id uuid.UUID) error
}

