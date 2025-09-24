package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type RequestStatus string

const (
	RequestStatusPendiente  RequestStatus = "PENDIENTE"
	RequestStatusAsignada   RequestStatus = "ASIGNADA"
	RequestStatusEnRuta     RequestStatus = "EN_RUTA"
	RequestStatusCompletada RequestStatus = "COMPLETADA"
	RequestStatusCancelada  RequestStatus = "CANCELADA"
)

type Language string

const (
	LanguageSpanish    Language = "es"
	LanguageEnglish    Language = "en"
	LanguagePortuguese Language = "pt"
	LanguageFrench     Language = "fr"
)

type Address struct {
	Street string `json:"street"`
	Number string `json:"number,omitempty"`
	City   string `json:"city"`
	Region string `json:"region"`
}

type Request struct {
	ID               uuid.UUID     `json:"id"`
	HotelID          *uuid.UUID    `json:"hotel_id,omitempty"`
	CompanyID        *uuid.UUID    `json:"company_id,omitempty"`
	Fecha            time.Time     `json:"fecha"`
	Origin           Address       `json:"origin"`
	Destination      Address       `json:"destination"`
	Pax              int           `json:"pax"`
	VehicleType      VehicleType   `json:"vehicle_type"`
	Language         *Language     `json:"language,omitempty"`
	Status           RequestStatus `json:"status"`
	AssignedDriverID *string       `json:"assigned_driver_id,omitempty"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`

	// Related data
	AssignedDriver *Driver         `json:"assigned_driver,omitempty"`
	Hotel          *Hotel          `json:"hotel,omitempty"`
	Company        *Company        `json:"company,omitempty"`
	Timeline       []TimelineEvent `json:"timeline,omitempty"`
}

type CreateRequestRequest struct {
	HotelID     *uuid.UUID  `json:"hotel_id,omitempty"`
	CompanyID   *uuid.UUID  `json:"company_id,omitempty"`
	Fecha       time.Time   `json:"fecha" validate:"required"`
	Origin      Address     `json:"origin" validate:"required"`
	Destination Address     `json:"destination" validate:"required"`
	Pax         int         `json:"pax" validate:"required,min=1"`
	VehicleType VehicleType `json:"vehicle_type" validate:"required"`
	Language    *Language   `json:"language,omitempty"`
}

type UpdateRequestRequest struct {
	Fecha       *time.Time   `json:"fecha,omitempty"`
	Origin      *Address     `json:"origin,omitempty"`
	Destination *Address     `json:"destination,omitempty"`
	Pax         *int         `json:"pax,omitempty" validate:"omitempty,min=1"`
	VehicleType *VehicleType `json:"vehicle_type,omitempty"`
	Language    *Language    `json:"language,omitempty"`
}

type AssignDriverRequest struct {
	DriverID string `json:"driver_id" validate:"required"`
}

type ChangeStatusRequest struct {
	NewStatus RequestStatus `json:"new_status" validate:"required"`
	Notes     *string       `json:"notes,omitempty"`
}

type ListRequestsRequest struct {
	Query       *string        `json:"query,omitempty"`
	Status      *RequestStatus `json:"status,omitempty"`
	VehicleType *VehicleType   `json:"vehicle_type,omitempty"`
	Page        int            `json:"page" validate:"min=1"`
	PageSize    int            `json:"page_size" validate:"min=1,max=100"`
	Sort        string         `json:"sort" validate:"omitempty,oneof=fecha created_at"`
}

// Helper methods for JSON marshaling/unmarshaling
func (r *Request) MarshalOrigin() ([]byte, error) {
	return json.Marshal(r.Origin)
}

func (r *Request) UnmarshalOrigin(data []byte) error {
	return json.Unmarshal(data, &r.Origin)
}

func (r *Request) MarshalDestination() ([]byte, error) {
	return json.Marshal(r.Destination)
}

func (r *Request) UnmarshalDestination(data []byte) error {
	return json.Unmarshal(data, &r.Destination)
}

// Validation methods
func (r *Request) CanTransitionTo(newStatus RequestStatus) bool {
	switch r.Status {
	case RequestStatusPendiente:
		return newStatus == RequestStatusAsignada || newStatus == RequestStatusCancelada
	case RequestStatusAsignada:
		return newStatus == RequestStatusEnRuta || newStatus == RequestStatusCancelada
	case RequestStatusEnRuta:
		return newStatus == RequestStatusCompletada || newStatus == RequestStatusCancelada
	case RequestStatusCompletada, RequestStatusCancelada:
		return false // Terminal states
	}
	return false
}

type RequestRepository interface {
	Create(request *Request) error
	GetByID(id uuid.UUID) (*Request, error)
	List(req ListRequestsRequest) ([]*Request, int, error)
	Update(id uuid.UUID, req UpdateRequestRequest) (*Request, error)
	Delete(id uuid.UUID) error
	AssignDriver(id uuid.UUID, driverID string) error
	ChangeStatus(id uuid.UUID, newStatus RequestStatus) error
	GetTimeline(id uuid.UUID) ([]TimelineEvent, error)
	AddTimelineEvent(id uuid.UUID, event TimelineEvent) error
}

