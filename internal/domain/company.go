package domain

import (
	"time"

	"github.com/google/uuid"
)

type CompanyStatus string

const (
	CompanyStatusActive    CompanyStatus = "ACTIVE"
	CompanyStatusSuspended CompanyStatus = "SUSPENDED"
)

type CompanySector string

const (
	CompanySectorHotel   CompanySector = "HOTEL"
	CompanySectorMineria CompanySector = "MINERIA"
	CompanySectorTurismo CompanySector = "TURISMO"
)

type Company struct {
	ID           uuid.UUID     `json:"id"`
	Name         string        `json:"name"`
	RUT          string        `json:"rut"`
	ContactEmail string        `json:"contact_email"`
	Status       CompanyStatus `json:"status"`
	Sector       CompanySector `json:"sector"`
	CreatedAt    time.Time     `json:"created_at"`
	UpdatedAt    time.Time     `json:"updated_at"`
}

type CreateCompanyRequest struct {
	Name         string        `json:"name" validate:"required,min=2,max=255"`
	RUT          string        `json:"rut" validate:"required,min=8,max=50"`
	ContactEmail string        `json:"contact_email" validate:"required,email"`
	Status       CompanyStatus `json:"status" validate:"required"`
	Sector       CompanySector `json:"sector" validate:"required"`
}

type UpdateCompanyRequest struct {
	Name         *string        `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	RUT          *string        `json:"rut,omitempty" validate:"omitempty,min=8,max=50"`
	ContactEmail *string        `json:"contact_email,omitempty" validate:"omitempty,email"`
	Status       *CompanyStatus `json:"status,omitempty"`
	Sector       *CompanySector `json:"sector,omitempty"`
}

type ListCompaniesRequest struct {
	Query    *string        `json:"query,omitempty"`
	Status   *CompanyStatus `json:"status,omitempty"`
	Sector   *CompanySector `json:"sector,omitempty"`
	OrgID    *uuid.UUID     `json:"org_id,omitempty"`
	Page     int            `json:"page" validate:"min=1"`
	PageSize int            `json:"page_size" validate:"min=1,max=100"`
	Sort     string         `json:"sort" validate:"omitempty,oneof=name rut created_at"`
}

type CompanyRepository interface {
	Create(company *Company) error
	GetByID(id uuid.UUID) (*Company, error)
	GetByRUT(rut string) (*Company, error)
	List(req ListCompaniesRequest) ([]*Company, int, error)
	Update(id uuid.UUID, req UpdateCompanyRequest) (*Company, error)
	Delete(id uuid.UUID) error
}
