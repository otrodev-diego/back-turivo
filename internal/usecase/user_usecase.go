package usecase

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type UserUseCase struct {
	userRepo              domain.UserRepository
	registrationTokenRepo domain.RegistrationTokenRepository
	passwordService       domain.PasswordService
	emailService          domain.EmailService
	logger                *zap.Logger
}

func NewUserUseCase(
	userRepo domain.UserRepository,
	registrationTokenRepo domain.RegistrationTokenRepository,
	passwordService domain.PasswordService,
	emailService domain.EmailService,
	logger *zap.Logger,
) *UserUseCase {
	return &UserUseCase{
		userRepo:              userRepo,
		registrationTokenRepo: registrationTokenRepo,
		passwordService:       passwordService,
		emailService:          emailService,
		logger:                logger,
	}
}

func (uc *UserUseCase) CreateUser(req domain.CreateUserRequest) (*domain.User, error) {
	uc.logger.Info("Creating user", zap.String("email", req.Email))

	// Check if user already exists
	existingUser, err := uc.userRepo.GetByEmail(req.Email)
	if err != nil && err != domain.ErrUserNotFound {
		uc.logger.Error("Failed to check existing user", zap.Error(err))
		return nil, domain.ErrInternalError
	}
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := uc.passwordService.HashPassword(req.Password)
	if err != nil {
		uc.logger.Error("Failed to hash password", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Create user
	user := &domain.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		Role:         req.Role,
		Status:       domain.UserStatusActive,
		OrgID:        req.OrgID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(user); err != nil {
		uc.logger.Error("Failed to create user", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	uc.logger.Info("User created successfully", zap.String("user_id", user.ID.String()))
	return user, nil
}

func (uc *UserUseCase) GetUserByID(id uuid.UUID) (*domain.User, error) {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrUserNotFound
		}
		uc.logger.Error("Failed to get user by ID", zap.Error(err), zap.String("user_id", id.String()))
		return nil, domain.ErrInternalError
	}

	return user, nil
}

func (uc *UserUseCase) GetUserByEmail(email string) (*domain.User, error) {
	user, err := uc.userRepo.GetByEmail(email)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrUserNotFound
		}
		uc.logger.Error("Failed to get user by email", zap.Error(err), zap.String("email", email))
		return nil, domain.ErrInternalError
	}

	return user, nil
}

func (uc *UserUseCase) ListUsers(req domain.ListUsersRequest) ([]*domain.User, int, error) {
	// Set default pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}
	if req.Sort == "" {
		req.Sort = "created_at"
	}

	users, total, err := uc.userRepo.List(req)
	if err != nil {
		uc.logger.Error("Failed to list users", zap.Error(err))
		return nil, 0, domain.ErrInternalError
	}

	return users, total, nil
}

func (uc *UserUseCase) UpdateUser(id uuid.UUID, req domain.UpdateUserRequest) (*domain.User, error) {
	uc.logger.Info("Updating user", zap.String("user_id", id.String()))

	// Check if user exists
	existingUser, err := uc.userRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrUserNotFound
		}
		uc.logger.Error("Failed to get user for update", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Check if email is being changed and if it already exists
	if req.Email != nil && *req.Email != existingUser.Email {
		emailUser, err := uc.userRepo.GetByEmail(*req.Email)
		if err != nil && err != domain.ErrUserNotFound {
			uc.logger.Error("Failed to check email uniqueness", zap.Error(err))
			return nil, domain.ErrInternalError
		}
		if emailUser != nil {
			return nil, domain.ErrUserAlreadyExists
		}
	}

	// Hash password if provided
	if req.Password != nil {
		hashedPassword, err := uc.passwordService.HashPassword(*req.Password)
		if err != nil {
			uc.logger.Error("Failed to hash password", zap.Error(err))
			return nil, domain.ErrInternalError
		}
		req.Password = &hashedPassword
	}

	user, err := uc.userRepo.Update(id, req)
	if err != nil {
		uc.logger.Error("Failed to update user", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	uc.logger.Info("User updated successfully", zap.String("user_id", user.ID.String()))
	return user, nil
}

func (uc *UserUseCase) DeleteUser(id uuid.UUID) error {
	uc.logger.Info("Deleting user", zap.String("user_id", id.String()))

	// Check if user exists
	_, err := uc.userRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrUserNotFound
		}
		uc.logger.Error("Failed to get user for deletion", zap.Error(err))
		return domain.ErrInternalError
	}

	if err := uc.userRepo.Delete(id); err != nil {
		uc.logger.Error("Failed to delete user", zap.Error(err))
		return domain.ErrInternalError
	}

	uc.logger.Info("User deleted successfully", zap.String("user_id", id.String()))
	return nil
}

// CreateUserWithInvitation creates a user and sends a welcome email with registration token
func (uc *UserUseCase) CreateUserWithInvitation(email string, role domain.UserRole, orgID *uuid.UUID) error {
	uc.logger.Info("🔵 === CreateUserWithInvitation UseCase Started ===", 
		zap.String("email", email),
		zap.String("role", string(role)),
		zap.Bool("has_org_id", orgID != nil),
	)

	// Check if user already exists
	uc.logger.Info("🔍 Checking if user already exists", zap.String("email", email))
	existingUser, err := uc.userRepo.GetByEmail(email)
	if err != nil && err != domain.ErrUserNotFound {
		uc.logger.Error("❌ FAILED to check existing user", zap.Error(err))
		return domain.ErrInternalError
	}
	if existingUser != nil {
		uc.logger.Warn("⚠️ User already exists", zap.String("email", email))
		return domain.ErrUserAlreadyExists
	}
	uc.logger.Info("✅ User doesn't exist, proceeding", zap.String("email", email))

	// Generate registration token
	uc.logger.Info("🔑 Generating registration token")
	token, err := uc.generateRegistrationToken()
	if err != nil {
		uc.logger.Error("❌ FAILED to generate registration token", zap.Error(err))
		return domain.ErrInternalError
	}
	uc.logger.Info("✅ Registration token generated", zap.String("token_length", fmt.Sprintf("%d", len(token))))

	// Create registration token record
	uc.logger.Info("💾 Creating registration token record in database")
	registrationToken := &domain.RegistrationToken{
		ID:        uuid.New(),
		Token:     token,
		Email:     email,
		OrgID:     orgID,
		Role:      role,
		ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // 24 hours expiration
		Used:      false,
		CreatedAt: time.Now().Unix(),
	}

	if err := uc.registrationTokenRepo.Create(registrationToken); err != nil {
		uc.logger.Error("❌ FAILED to create registration token in database", 
			zap.Error(err),
			zap.String("email", email),
			zap.String("token_id", registrationToken.ID.String()),
		)
		return domain.ErrInternalError
	}
	uc.logger.Info("✅ Registration token saved to database", 
		zap.String("token_id", registrationToken.ID.String()),
		zap.String("email", email),
	)

	// Send welcome email
	uc.logger.Info("📧 Sending welcome email", zap.String("email", email))
	name := email // Use email as name since we don't have the name yet
	if err := uc.emailService.SendWelcomeEmail(email, name, token); err != nil {
		uc.logger.Error("❌ FAILED to send welcome email", 
			zap.Error(err),
			zap.String("email", email),
		)
		// We don't return error here because the token is already created
		// The admin can resend the email manually if needed
	} else {
		uc.logger.Info("✅ Welcome email sent successfully", zap.String("email", email))
	}

	uc.logger.Info("🎉 User invitation process completed successfully", zap.String("email", email))
	return nil
}

// CompleteRegistration completes user registration using a token
func (uc *UserUseCase) CompleteRegistration(req domain.CompleteRegistrationRequest) (*domain.User, error) {
	uc.logger.Info("Completing user registration", zap.String("token", req.Token))

	// Get registration token
	regToken, err := uc.registrationTokenRepo.GetByToken(req.Token)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, fmt.Errorf("invalid or expired registration token")
		}
		uc.logger.Error("Failed to get registration token", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Check if token is still valid
	if regToken.Used {
		return nil, fmt.Errorf("registration token has already been used")
	}

	if time.Now().Unix() > regToken.ExpiresAt {
		return nil, fmt.Errorf("registration token has expired")
	}

	// Check if user already exists (shouldn't happen, but just in case)
	existingUser, err := uc.userRepo.GetByEmail(regToken.Email)
	if err != nil && err != domain.ErrUserNotFound {
		uc.logger.Error("Failed to check existing user", zap.Error(err))
		return nil, domain.ErrInternalError
	}
	if existingUser != nil {
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := uc.passwordService.HashPassword(req.Password)
	if err != nil {
		uc.logger.Error("Failed to hash password", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Create user
	user := &domain.User{
		ID:           uuid.New(),
		Name:         req.Name,
		Email:        regToken.Email,
		PasswordHash: hashedPassword,
		Role:         regToken.Role,
		Status:       domain.UserStatusActive,
		OrgID:        regToken.OrgID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := uc.userRepo.Create(user); err != nil {
		uc.logger.Error("Failed to create user", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Mark token as used
	if err := uc.registrationTokenRepo.MarkAsUsed(req.Token); err != nil {
		uc.logger.Error("Failed to mark registration token as used", zap.Error(err))
		// Don't return error here as the user is already created
	}

	uc.logger.Info("User registration completed successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("email", user.Email),
	)
	return user, nil
}

// ValidateRegistrationToken validates if a registration token is valid
func (uc *UserUseCase) ValidateRegistrationToken(token string) (*domain.RegistrationToken, error) {
	regToken, err := uc.registrationTokenRepo.GetByToken(token)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, fmt.Errorf("invalid registration token")
		}
		uc.logger.Error("Failed to get registration token", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	if regToken.Used {
		return nil, fmt.Errorf("registration token has already been used")
	}

	if time.Now().Unix() > regToken.ExpiresAt {
		return nil, fmt.Errorf("registration token has expired")
	}

	return regToken, nil
}

// generateRegistrationToken generates a secure random token
func (uc *UserUseCase) generateRegistrationToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
