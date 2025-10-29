package domain

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	User         *User  `json:"user"`
}

type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

type RefreshTokenResponse struct {
	AccessToken string `json:"access_token"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=8"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}

type JWTClaims struct {
	UserID         uuid.UUID       `json:"sub"`
	Role           UserRole        `json:"role"`
	OrgID          *uuid.UUID      `json:"org_id,omitempty"`
	CompanyProfile *CompanyProfile `json:"company_profile,omitempty"`
	Exp            int64           `json:"exp"`
	Iat            int64           `json:"iat"`
}

// Implement jwt.Claims interface
func (c JWTClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.Exp, 0)), nil
}

func (c JWTClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	return jwt.NewNumericDate(time.Unix(c.Iat, 0)), nil
}

func (c JWTClaims) GetNotBefore() (*jwt.NumericDate, error) {
	return nil, nil
}

func (c JWTClaims) GetIssuer() (string, error) {
	return "", nil
}

func (c JWTClaims) GetSubject() (string, error) {
	return c.UserID.String(), nil
}

func (c JWTClaims) GetAudience() (jwt.ClaimStrings, error) {
	return nil, nil
}

type RefreshToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type PasswordResetToken struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
	Used      bool      `json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

type AuthService interface {
	Login(req LoginRequest) (*LoginResponse, error)
	RefreshToken(req RefreshTokenRequest) (*RefreshTokenResponse, error)
	Logout(refreshToken string) error
	ValidateAccessToken(token string) (*JWTClaims, error)
	GenerateAccessToken(user *User) (string, error)
	GenerateRefreshToken(userID uuid.UUID) (*RefreshToken, error)
	ForgotPassword(req ForgotPasswordRequest) (*ForgotPasswordResponse, error)
	ResetPassword(req ResetPasswordRequest) (*ResetPasswordResponse, error)
}

type RefreshTokenRepository interface {
	Create(token *RefreshToken) error
	GetByToken(token string) (*RefreshToken, error)
	Delete(token string) error
	DeleteByUserID(userID uuid.UUID) error
}

type PasswordResetTokenRepository interface {
	Create(token *PasswordResetToken) error
	GetByToken(token string) (*PasswordResetToken, error)
	MarkAsUsed(token string) error
	DeleteExpiredTokens() error
}

type PasswordService interface {
	HashPassword(password string) (string, error)
	VerifyPassword(hashedPassword, password string) error
}
