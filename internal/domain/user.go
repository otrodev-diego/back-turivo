package domain

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	UserRoleAdmin   UserRole = "ADMIN"
	UserRoleUser    UserRole = "USER"
	UserRoleDriver  UserRole = "DRIVER"
	UserRoleCompany UserRole = "COMPANY"
)

type UserStatus string

const (
	UserStatusActive  UserStatus = "ACTIVE"
	UserStatusBlocked UserStatus = "BLOCKED"
)

type CompanyProfile string

const (
	CompanyProfileAdmin CompanyProfile = "COMPANY_ADMIN"
	CompanyProfileUser  CompanyProfile = "COMPANY_USER"
)

type User struct {
	ID             uuid.UUID       `json:"id"`
	Name           string          `json:"name"`
	Email          string          `json:"email"`
	PasswordHash   string          `json:"-"`
	Role           UserRole        `json:"role"`
	Status         UserStatus      `json:"status"`
	OrgID          *uuid.UUID      `json:"org_id,omitempty"`
	CompanyProfile *CompanyProfile `json:"company_profile,omitempty"`
	CreatedAt      time.Time       `json:"created_at"`
	UpdatedAt      time.Time       `json:"updated_at"`
}

type CreateUserRequest struct {
	Name           string          `json:"name" validate:"required,min=2,max=255"`
	Email          string          `json:"email" validate:"required,email"`
	Password       string          `json:"password" validate:"required,min=8"`
	Role           UserRole        `json:"role" validate:"required"`
	OrgID          *uuid.UUID      `json:"org_id,omitempty"`
	CompanyProfile *CompanyProfile `json:"company_profile,omitempty"`
}

type UpdateUserRequest struct {
	Name           *string         `json:"name,omitempty" validate:"omitempty,min=2,max=255"`
	Email          *string         `json:"email,omitempty" validate:"omitempty,email"`
	Password       *string         `json:"password,omitempty" validate:"omitempty,min=8"`
	Role           *UserRole       `json:"role,omitempty"`
	Status         *UserStatus     `json:"status,omitempty"`
	OrgID          *uuid.UUID      `json:"org_id,omitempty"`
	CompanyProfile *CompanyProfile `json:"company_profile,omitempty"`
}

type ListUsersRequest struct {
	Query    *string     `json:"query,omitempty"`
	Role     *UserRole   `json:"role,omitempty"`
	Status   *UserStatus `json:"status,omitempty"`
	OrgID    *uuid.UUID  `json:"org_id,omitempty"`
	Page     int         `json:"page" validate:"min=1"`
	PageSize int         `json:"page_size" validate:"min=1,max=100"`
	Sort     string      `json:"sort" validate:"omitempty,oneof=name email created_at"`
}

type UserRepository interface {
	Create(user *User) error
	GetByID(id uuid.UUID) (*User, error)
	GetByEmail(email string) (*User, error)
	List(req ListUsersRequest) ([]*User, int, error)
	Update(id uuid.UUID, req UpdateUserRequest) (*User, error)
	UpdateUser(user *User) error
	Delete(id uuid.UUID) error
}
