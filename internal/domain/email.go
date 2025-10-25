package domain

import "github.com/google/uuid"

type EmailService interface {
	SendWelcomeEmail(email, name string, registrationToken string) error
	SendReservationCreated(to string, reservation *Reservation, user *User) error
	SendReservationNotification(to string, reservation *Reservation, user *User) error
	SendSupportRequest(to string, request *SupportRequest, user *User) error
}

type WelcomeEmailData struct {
	Name              string
	Email             string
	RegistrationToken string
	RegistrationURL   string
}

type ReservationEmailData struct {
	UserName       string
	UserEmail      string
	ReservationID  string
	Pickup         string
	Destination    string
	DateTime       string
	Passengers     int
	VehicleType    string
	Amount         *float64
	Status         string
	Notes          string
	Stops          int
	HasSpecialLang bool
}

type SupportEmailData struct {
	UserID      string
	UserName    string
	UserEmail   string
	Descripcion string
	Detalle     string
	Timestamp   string
}

type SupportRequest struct {
	UserID      string `json:"user_id"`
	Descripcion string `json:"descripcion" validate:"required,min=5,max=255"`
	Detalle     string `json:"detalle" validate:"required,min=10,max=2000"`
}

// Registration token domain
type RegistrationToken struct {
	ID             uuid.UUID       `json:"id"`
	Token          string          `json:"token"`
	Email          string          `json:"email"`
	OrgID          *uuid.UUID      `json:"org_id,omitempty"`
	Role           UserRole        `json:"role"`
	CompanyProfile *CompanyProfile `json:"company_profile,omitempty"`
	ExpiresAt      int64           `json:"expires_at"`
	Used           bool            `json:"used"`
	CreatedAt      int64           `json:"created_at"`
}

type RegistrationTokenRepository interface {
	Create(token *RegistrationToken) error
	GetByToken(token string) (*RegistrationToken, error)
	MarkAsUsed(token string) error
	DeleteExpired() error
	ListAll() ([]*RegistrationToken, error)
}

type CompleteRegistrationRequest struct {
	Token    string `json:"token" validate:"required"`
	Name     string `json:"name" validate:"required,min=2,max=255"`
	Password string `json:"password" validate:"required,min=8"`
}
