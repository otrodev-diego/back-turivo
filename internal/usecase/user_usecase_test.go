package usecase

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

// Mock implementations
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uuid.UUID) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) List(req domain.ListUsersRequest) ([]*domain.User, int, error) {
	args := m.Called(req)
	return args.Get(0).([]*domain.User), args.Int(1), args.Error(2)
}

func (m *MockUserRepository) Update(id uuid.UUID, req domain.UpdateUserRequest) (*domain.User, error) {
	args := m.Called(id, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Delete(id uuid.UUID) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockPasswordService struct {
	mock.Mock
}

func (m *MockPasswordService) HashPassword(password string) (string, error) {
	args := m.Called(password)
	return args.String(0), args.Error(1)
}

func (m *MockPasswordService) VerifyPassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func TestUserUseCase_CreateUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPasswordService := new(MockPasswordService)
	logger := zap.NewNop()

	useCase := NewUserUseCase(mockRepo, mockPasswordService, logger)

	t.Run("should create user successfully", func(t *testing.T) {
		req := domain.CreateUserRequest{
			Name:     "Test User",
			Email:    "test@example.com",
			Password: "password123",
			Role:     domain.UserRoleUser,
		}

		// Mock expectations
		mockRepo.On("GetByEmail", req.Email).Return(nil, domain.ErrUserNotFound)
		mockPasswordService.On("HashPassword", req.Password).Return("hashedpassword", nil)
		mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

		// Execute
		user, err := useCase.CreateUser(req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, user)
		assert.Equal(t, req.Name, user.Name)
		assert.Equal(t, req.Email, user.Email)
		assert.Equal(t, req.Role, user.Role)
		assert.Equal(t, domain.UserStatusActive, user.Status)
		assert.Equal(t, "hashedpassword", user.PasswordHash)

		mockRepo.AssertExpectations(t)
		mockPasswordService.AssertExpectations(t)
	})

	t.Run("should return error when user already exists", func(t *testing.T) {
		req := domain.CreateUserRequest{
			Name:     "Test User",
			Email:    "existing@example.com",
			Password: "password123",
			Role:     domain.UserRoleUser,
		}

		existingUser := &domain.User{
			ID:    uuid.New(),
			Email: req.Email,
		}

		// Mock expectations
		mockRepo.On("GetByEmail", req.Email).Return(existingUser, nil)

		// Execute
		user, err := useCase.CreateUser(req)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.ErrUserAlreadyExists, err)
		assert.Nil(t, user)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_GetUserByID(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPasswordService := new(MockPasswordService)
	logger := zap.NewNop()

	useCase := NewUserUseCase(mockRepo, mockPasswordService, logger)

	t.Run("should get user successfully", func(t *testing.T) {
		userID := uuid.New()
		expectedUser := &domain.User{
			ID:        userID,
			Name:      "Test User",
			Email:     "test@example.com",
			Role:      domain.UserRoleUser,
			Status:    domain.UserStatusActive,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		// Mock expectations
		mockRepo.On("GetByID", userID).Return(expectedUser, nil)

		// Execute
		user, err := useCase.GetUserByID(userID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should return error when user not found", func(t *testing.T) {
		userID := uuid.New()

		// Mock expectations
		mockRepo.On("GetByID", userID).Return(nil, domain.ErrNotFound)

		// Execute
		user, err := useCase.GetUserByID(userID)

		// Assert
		assert.Error(t, err)
		assert.Equal(t, domain.ErrUserNotFound, err)
		assert.Nil(t, user)

		mockRepo.AssertExpectations(t)
	})
}

func TestUserUseCase_ListUsers(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockPasswordService := new(MockPasswordService)
	logger := zap.NewNop()

	useCase := NewUserUseCase(mockRepo, mockPasswordService, logger)

	t.Run("should list users successfully", func(t *testing.T) {
		req := domain.ListUsersRequest{
			Page:     1,
			PageSize: 20,
			Sort:     "created_at",
		}

		expectedUsers := []*domain.User{
			{
				ID:    uuid.New(),
				Name:  "User 1",
				Email: "user1@example.com",
			},
			{
				ID:    uuid.New(),
				Name:  "User 2",
				Email: "user2@example.com",
			},
		}

		// Mock expectations
		mockRepo.On("List", mock.MatchedBy(func(r domain.ListUsersRequest) bool {
			return r.Page == 1 && r.PageSize == 20 && r.Sort == "created_at"
		})).Return(expectedUsers, 2, nil)

		// Execute
		users, total, err := useCase.ListUsers(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Equal(t, 2, total)

		mockRepo.AssertExpectations(t)
	})

	t.Run("should set defaults for pagination", func(t *testing.T) {
		// Create fresh mocks for this test
		mockRepo := new(MockUserRepository)
		mockPasswordService := new(MockPasswordService)
		useCase := NewUserUseCase(mockRepo, mockPasswordService, logger)

		req := domain.ListUsersRequest{
			Page:     0, // Invalid
			PageSize: 0, // Invalid
		}

		expectedUsers := []*domain.User{}

		// Mock expectations - should be called with corrected values
		mockRepo.On("List", mock.MatchedBy(func(r domain.ListUsersRequest) bool {
			return r.Page == 1 && r.PageSize == 20 && r.Sort == "created_at"
		})).Return(expectedUsers, 0, nil)

		// Execute
		users, total, err := useCase.ListUsers(req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, expectedUsers, users)
		assert.Equal(t, 0, total)

		mockRepo.AssertExpectations(t)
	})
}
