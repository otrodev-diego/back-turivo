package domain

import (
	"time"

	"github.com/google/uuid"
)

type Hotel struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	City         string    `json:"city"`
	ContactEmail string    `json:"contact_email"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type CreateHotelRequest struct {
	Name         string `json:"name" validate:"required,min=2,max=255"`
	City         string `json:"city" validate:"required,min=2,max=255"`
	ContactEmail string `json:"contact_email" validate:"required,email"`
}

type UpdateHotelRequest struct {
	Name         *string `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	City         *string `json:"city,omitempty" validate:"omitempty,min=2,max=255"`
	ContactEmail *string `json:"contact_email,omitempty" validate:"omitempty,email"`
}

type ListHotelsRequest struct {
	Query    *string `json:"query,omitempty"`
	Page     int     `json:"page" validate:"min=1"`
	PageSize int     `json:"page_size" validate:"min=1,max=100"`
	Sort     string  `json:"sort" validate:"omitempty,oneof=name city created_at"`
}

type HotelRepository interface {
	Create(hotel *Hotel) error
	GetByID(id uuid.UUID) (*Hotel, error)
	List(req ListHotelsRequest) ([]*Hotel, int, error)
	Update(id uuid.UUID, req UpdateHotelRequest) (*Hotel, error)
	Delete(id uuid.UUID) error
}
