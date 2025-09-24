package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type SupportHandler struct {
	emailService domain.EmailService
	userRepo     domain.UserRepository
	validator    *validator.Validate
	logger       *zap.Logger
}

func NewSupportHandler(emailService domain.EmailService, userRepo domain.UserRepository, validator *validator.Validate, logger *zap.Logger) *SupportHandler {
	return &SupportHandler{
		emailService: emailService,
		userRepo:     userRepo,
		validator:    validator,
		logger:       logger,
	}
}

// ContactSupport godoc
// @Summary Send support request
// @Description Send a support request email to operations team
// @Tags support
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.SupportRequest true "Support request data"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/support/contact [post]
func (h *SupportHandler) ContactSupport(c *gin.Context) {
	// Get user ID from JWT token (set by auth middleware)
	userID, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("User ID not found in context")
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User not authenticated",
		})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		h.logger.Warn("Invalid user ID type in context")
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	var req domain.SupportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for support request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Set user ID from context
	req.UserID = userUUID.String()

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for support request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	// Get user information
	user, err := h.userRepo.GetByID(userUUID)
	if err != nil {
		h.logger.Error("Failed to get user information", zap.Error(err), zap.String("user_id", userUUID.String()))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	// Send support request email
	operationsEmail := "djaramontenegro@gmail.com"
	if err := h.emailService.SendSupportRequest(operationsEmail, &req, user); err != nil {
		h.logger.Error("Failed to send support request email", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to send support request",
		})
		return
	}

	h.logger.Info("Support request sent successfully",
		zap.String("user_id", userUUID.String()),
		zap.String("user_email", user.Email),
		zap.String("descripcion", req.Descripcion),
	)

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Tu solicitud de soporte ha sido enviada exitosamente",
	})
}
