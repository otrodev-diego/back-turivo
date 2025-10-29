package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
	"turivo-backend/internal/usecase"
)

type AuthHandler struct {
	authUseCase *usecase.AuthUseCase
	validator   *validator.Validate
	logger      *zap.Logger
}

func NewAuthHandler(authUseCase *usecase.AuthUseCase, validator *validator.Validate, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authUseCase: authUseCase,
		validator:   validator,
		logger:      logger,
	}
}

// Login godoc
// @Summary User login
// @Description Authenticate user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.LoginRequest true "Login credentials"
// @Success 200 {object} domain.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for login", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for login", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	response, err := h.authUseCase.Login(req)
	if err != nil {
		switch err {
		case domain.ErrInvalidCredentials:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Invalid credentials",
			})
		case domain.ErrUserBlocked:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "User account is blocked",
			})
		default:
			h.logger.Error("Login failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken godoc
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} domain.RefreshTokenResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req domain.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for refresh token", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for refresh token", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	response, err := h.authUseCase.RefreshToken(req)
	if err != nil {
		switch err {
		case domain.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Invalid or expired refresh token",
			})
		case domain.ErrUserBlocked:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "User account is blocked",
			})
		default:
			h.logger.Error("Token refresh failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout godoc
// @Summary User logout
// @Description Invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.RefreshTokenRequest true "Refresh token to invalidate"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req domain.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for logout", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for logout", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	if err := h.authUseCase.Logout(req.RefreshToken); err != nil {
		h.logger.Error("Logout failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Logged out successfully",
	})
}

// ForgotPassword godoc
// @Summary Request password reset
// @Description Send password reset email to user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.ForgotPasswordRequest true "Email address"
// @Success 200 {object} domain.ForgotPasswordResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/forgot-password [post]
func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req domain.ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for forgot password", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for forgot password", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	response, err := h.authUseCase.ForgotPassword(req)
	if err != nil {
		h.logger.Error("Forgot password failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, response)
}

// ResetPassword godoc
// @Summary Reset password with token
// @Description Reset user password using reset token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body domain.ResetPasswordRequest true "Reset password data"
// @Success 200 {object} domain.ResetPasswordResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/auth/reset-password [post]
func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req domain.ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for reset password", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for reset password", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	response, err := h.authUseCase.ResetPassword(req)
	if err != nil {
		switch err {
		case domain.ErrUnauthorized:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "Invalid or expired reset token",
			})
		case domain.ErrUserBlocked:
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error: "User account is blocked",
			})
		default:
			h.logger.Error("Reset password failed", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, response)
}

// Common response structures
type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}
