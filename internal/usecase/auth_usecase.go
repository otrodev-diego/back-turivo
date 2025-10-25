package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type AuthUseCase struct {
	userRepo         domain.UserRepository
	refreshTokenRepo domain.RefreshTokenRepository
	passwordService  domain.PasswordService
	jwtSecret        string
	accessTokenTTL   time.Duration
	refreshTokenTTL  time.Duration
	logger           *zap.Logger
}

func NewAuthUseCase(
	userRepo domain.UserRepository,
	refreshTokenRepo domain.RefreshTokenRepository,
	passwordService domain.PasswordService,
	jwtSecret string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	logger *zap.Logger,
) *AuthUseCase {
	return &AuthUseCase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		passwordService:  passwordService,
		jwtSecret:        jwtSecret,
		accessTokenTTL:   accessTokenTTL,
		refreshTokenTTL:  refreshTokenTTL,
		logger:           logger,
	}
}

func (uc *AuthUseCase) Login(req domain.LoginRequest) (*domain.LoginResponse, error) {
	uc.logger.Info("User login attempt", zap.String("email", req.Email))

	// Get user by email
	user, err := uc.userRepo.GetByEmail(req.Email)
	if err != nil {
		if err == domain.ErrNotFound {
			uc.logger.Warn("Login attempt with non-existent email", zap.String("email", req.Email))
			return nil, domain.ErrInvalidCredentials
		}
		uc.logger.Error("Failed to get user for login", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Check if user is blocked
	if user.Status == domain.UserStatusBlocked {
		uc.logger.Warn("Login attempt by blocked user", zap.String("user_id", user.ID.String()))
		return nil, domain.ErrUserBlocked
	}

	// Verify password
	if err := uc.passwordService.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		uc.logger.Warn("Invalid password attempt", zap.String("email", req.Email))
		return nil, domain.ErrInvalidCredentials
	}

	// Generate access token
	accessToken, err := uc.GenerateAccessToken(user)
	if err != nil {
		uc.logger.Error("Failed to generate access token", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Generate refresh token
	refreshToken, err := uc.GenerateRefreshToken(user.ID)
	if err != nil {
		uc.logger.Error("Failed to generate refresh token", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	uc.logger.Info("User logged in successfully", zap.String("user_id", user.ID.String()))

	return &domain.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken.Token,
		User:         user,
	}, nil
}

func (uc *AuthUseCase) RefreshToken(req domain.RefreshTokenRequest) (*domain.RefreshTokenResponse, error) {
	uc.logger.Info("Refresh token request")

	// Get refresh token from repository
	refreshToken, err := uc.refreshTokenRepo.GetByToken(req.RefreshToken)
	if err != nil {
		if err == domain.ErrNotFound {
			uc.logger.Warn("Invalid refresh token used")
			return nil, domain.ErrUnauthorized
		}
		uc.logger.Error("Failed to get refresh token", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Check if token is expired
	if refreshToken.ExpiresAt.Before(time.Now()) {
		uc.logger.Warn("Expired refresh token used", zap.String("user_id", refreshToken.UserID.String()))
		// Delete expired token
		_ = uc.refreshTokenRepo.Delete(refreshToken.Token)
		return nil, domain.ErrUnauthorized
	}

	// Get user
	user, err := uc.userRepo.GetByID(refreshToken.UserID)
	if err != nil {
		if err == domain.ErrNotFound {
			uc.logger.Warn("Refresh token for non-existent user", zap.String("user_id", refreshToken.UserID.String()))
			return nil, domain.ErrUnauthorized
		}
		uc.logger.Error("Failed to get user for refresh", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Check if user is blocked
	if user.Status == domain.UserStatusBlocked {
		uc.logger.Warn("Refresh token attempt by blocked user", zap.String("user_id", user.ID.String()))
		return nil, domain.ErrUserBlocked
	}

	// Generate new access token
	accessToken, err := uc.GenerateAccessToken(user)
	if err != nil {
		uc.logger.Error("Failed to generate access token", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	uc.logger.Info("Token refreshed successfully", zap.String("user_id", user.ID.String()))

	return &domain.RefreshTokenResponse{
		AccessToken: accessToken,
	}, nil
}

func (uc *AuthUseCase) Logout(refreshToken string) error {
	uc.logger.Info("User logout")

	if err := uc.refreshTokenRepo.Delete(refreshToken); err != nil {
		uc.logger.Error("Failed to delete refresh token", zap.Error(err))
		return domain.ErrInternalError
	}

	return nil
}

func (uc *AuthUseCase) ValidateAccessToken(tokenString string) (*domain.JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &domain.JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrUnauthorized
		}
		return []byte(uc.jwtSecret), nil
	})

	if err != nil {
		return nil, domain.ErrUnauthorized
	}

	if claims, ok := token.Claims.(*domain.JWTClaims); ok && token.Valid {
		// Check if token is expired
		if time.Unix(claims.Exp, 0).Before(time.Now()) {
			return nil, domain.ErrUnauthorized
		}
		return claims, nil
	}

	return nil, domain.ErrUnauthorized
}

func (uc *AuthUseCase) GenerateAccessToken(user *domain.User) (string, error) {
	now := time.Now()
	uc.logger.Info("Generating access token", zap.Duration("accessTokenTTL", uc.accessTokenTTL))
	claims := &domain.JWTClaims{
		UserID:         user.ID,
		Role:           user.Role,
		OrgID:          user.OrgID,
		CompanyProfile: user.CompanyProfile,
		Exp:            now.Add(uc.accessTokenTTL).Unix(),
		Iat:            now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(uc.jwtSecret))
}

func (uc *AuthUseCase) GenerateRefreshToken(userID uuid.UUID) (*domain.RefreshToken, error) {
	// Generate random token
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return nil, err
	}
	tokenString := hex.EncodeToString(bytes)

	refreshToken := &domain.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(uc.refreshTokenTTL),
		CreatedAt: time.Now(),
	}

	if err := uc.refreshTokenRepo.Create(refreshToken); err != nil {
		return nil, err
	}

	return refreshToken, nil
}
