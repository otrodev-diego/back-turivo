package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
	"turivo-backend/internal/interface/http/middleware"
	"turivo-backend/internal/usecase"
)

type UserHandler struct {
	userUseCase *usecase.UserUseCase
	validator   *validator.Validate
	logger      *zap.Logger
}

func NewUserHandler(userUseCase *usecase.UserUseCase, validator *validator.Validate, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		validator:   validator,
		logger:      logger,
	}
}

// ListUsers godoc
// @Summary List users
// @Description Get paginated list of users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string false "Search query"
// @Param role query string false "Filter by role" Enums(ADMIN,USER,DRIVER,COMPANY)
// @Param status query string false "Filter by status" Enums(ACTIVE,BLOCKED)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param sort query string false "Sort by field" Enums(name,email,created_at)
// @Success 200 {object} PaginatedResponse{data=[]domain.User}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse query parameters
	req := domain.ListUsersRequest{
		Page:     1,
		PageSize: 20,
		Sort:     "created_at",
	}

	if q := c.Query("q"); q != "" {
		req.Query = &q
	}

	if role := c.Query("role"); role != "" {
		userRole := domain.UserRole(role)
		req.Role = &userRole
	}

	if status := c.Query("status"); status != "" {
		userStatus := domain.UserStatus(status)
		req.Status = &userStatus
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			req.Page = p
		}
	}

	if pageSize := c.Query("page_size"); pageSize != "" {
		if ps, err := strconv.Atoi(pageSize); err == nil && ps > 0 && ps <= 100 {
			req.PageSize = ps
		}
	}

	if sort := c.Query("sort"); sort != "" {
		req.Sort = sort
	}

	// Get organization scope from middleware
	if orgID, exists := c.Get("scope_org_id"); exists && orgID != nil {
		req.OrgID = orgID.(*uuid.UUID)
		h.logger.Info("Filtering users by organization", zap.String("org_id", req.OrgID.String()))
	} else {
		h.logger.Info("No organization scope - showing all users (admin)")
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for list users", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	users, total, err := h.userUseCase.ListUsers(req)
	if err != nil {
		h.logger.Error("Failed to list users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	totalPages := (total + req.PageSize - 1) / req.PageSize

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       users,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// CreateUser godoc
// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateUserRequest true "User data"
// @Success 201 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req domain.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for create user", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for create user", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	user, err := h.userUseCase.CreateUser(req)
	if err != nil {
		switch err {
		case domain.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "User already exists",
			})
		default:
			h.logger.Error("Failed to create user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

// GetUser godoc
// @Summary Get user
// @Description Get user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	user, err := h.userUseCase.GetUserByID(id)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "User not found",
			})
		default:
			h.logger.Error("Failed to get user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update user
// @Description Update user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Param request body domain.UpdateUserRequest true "User data"
// @Success 200 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id} [patch]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	var req domain.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for update user", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for update user", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	user, err := h.userUseCase.UpdateUser(id, req)
	if err != nil {
		switch err {
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "User not found",
			})
		case domain.ErrUserAlreadyExists:
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Email already exists",
			})
		default:
			h.logger.Error("Failed to update user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

// DeleteUser godoc
// @Summary Delete user
// @Description Delete user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid user ID",
		})
		return
	}

	// Check if user is trying to delete themselves
	userID, _ := middleware.GetUserID(c)
	if userID == id {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Cannot delete your own account",
		})
		return
	}

	if err := h.userUseCase.DeleteUser(id); err != nil {
		switch err {
		case domain.ErrUserNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "User not found",
			})
		default:
			h.logger.Error("Failed to delete user", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "User deleted successfully",
	})
}

// CreateUserInvitation godoc
// @Summary Create user invitation
// @Description Create a user invitation and send welcome email
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateUserInvitationRequest true "User invitation data"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/invite [post]
func (h *UserHandler) CreateUserInvitation(c *gin.Context) {
	h.logger.Info("=== CreateUserInvitation endpoint called ===")
	h.logger.Info("Request details",
		zap.String("method", c.Request.Method),
		zap.String("url", c.Request.URL.String()),
		zap.String("remote_addr", c.ClientIP()),
		zap.String("user_agent", c.GetHeader("User-Agent")),
	)

	var req CreateUserInvitationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("❌ FAILED to bind JSON request body",
			zap.Error(err),
			zap.String("raw_body", "check content-type and body format"),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	h.logger.Info("✅ Request body parsed successfully",
		zap.String("email", req.Email),
		zap.String("role", req.Role),
		zap.String("org_id", req.OrgID),
	)

	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("❌ VALIDATION FAILED",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("role", req.Role),
		)
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	h.logger.Info("✅ Validation passed")

	// Convert string role to domain.UserRole
	userRole := domain.UserRole(req.Role)
	h.logger.Info("🔄 Role converted", zap.String("role", string(userRole)))

	// Convert org_id string to UUID if provided
	var orgID *uuid.UUID
	if req.OrgID != "" {
		h.logger.Info("🔄 Parsing org_id", zap.String("org_id", req.OrgID))
		id, err := uuid.Parse(req.OrgID)
		if err != nil {
			h.logger.Error("❌ INVALID organization ID format",
				zap.Error(err),
				zap.String("org_id", req.OrgID),
			)
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid organization ID format",
			})
			return
		}
		orgID = &id
		h.logger.Info("✅ org_id parsed successfully", zap.String("org_id", id.String()))
	} else {
		h.logger.Info("ℹ️ No org_id provided")
	}

	h.logger.Info("🚀 About to call CreateUserWithInvitation",
		zap.String("email", req.Email),
		zap.String("role", string(userRole)),
		zap.Bool("has_org_id", orgID != nil),
	)

	// Vamos a probar con el usecase real ahora
	err := h.userUseCase.CreateUserWithInvitation(req.Email, userRole, orgID)
	if err != nil {
		switch err {
		case domain.ErrUserAlreadyExists:
			h.logger.Warn("⚠️ User already exists", zap.String("email", req.Email))
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "User already exists",
			})
		default:
			h.logger.Error("❌ FAILED to create user invitation",
				zap.Error(err),
				zap.String("email", req.Email),
				zap.String("error_type", fmt.Sprintf("%T", err)),
			)
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Internal server error",
				Details: err.Error(),
			})
		}
		return
	}

	h.logger.Info("🎉 User invitation created successfully", zap.String("email", req.Email))
	c.JSON(http.StatusOK, MessageResponse{
		Message: "User invitation sent successfully",
	})

	// Código original comentado temporalmente
	/*
		err := h.userUseCase.CreateUserWithInvitation(req.Email, userRole, orgID)
		if err != nil {
			switch err {
			case domain.ErrUserAlreadyExists:
				h.logger.Warn("User already exists", zap.String("email", req.Email))
				c.JSON(http.StatusConflict, ErrorResponse{
					Error: "User already exists",
				})
			default:
				h.logger.Error("Failed to create user invitation", zap.Error(err))
				c.JSON(http.StatusInternalServerError, ErrorResponse{
					Error:   "Internal server error",
					Details: err.Error(),
				})
			}
			return
		}

		h.logger.Info("User invitation created successfully", zap.String("email", req.Email))
		c.JSON(http.StatusOK, MessageResponse{
			Message: "User invitation sent successfully",
		})
	*/
}

// CompleteRegistration godoc
// @Summary Complete user registration
// @Description Complete user registration using a token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.CompleteRegistrationRequest true "Registration data"
// @Success 201 {object} domain.User
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/complete-registration [post]
func (h *UserHandler) CompleteRegistration(c *gin.Context) {
	var req domain.CompleteRegistrationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for complete registration", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for complete registration", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	user, err := h.userUseCase.CompleteRegistration(req)
	if err != nil {
		h.logger.Error("Failed to complete registration", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, user)
}

// ValidateRegistrationToken godoc
// @Summary Validate registration token
// @Description Validate if a registration token is valid
// @Tags auth
// @Accept json
// @Produce json
// @Param token query string true "Registration token"
// @Success 200 {object} domain.RegistrationToken
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/validate-token [get]
func (h *UserHandler) ValidateRegistrationToken(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Token parameter is required",
		})
		return
	}

	regToken, err := h.userUseCase.ValidateRegistrationToken(token)
	if err != nil {
		h.logger.Error("Failed to validate registration token", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, regToken)
}

// Request structures
type CreateUserInvitationRequest struct {
	Email string `json:"email" validate:"required,email"`
	Role  string `json:"role" validate:"required,oneof=ADMIN USER DRIVER COMPANY"`
	OrgID string `json:"org_id,omitempty"`
}

// Common response structures
type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
}
