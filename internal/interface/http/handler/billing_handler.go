package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
	"turivo-backend/internal/usecase"
)

type BillingHandler struct {
	paymentUseCase *usecase.PaymentUseCase
	logger         *zap.Logger
}

func NewBillingHandler(paymentUseCase *usecase.PaymentUseCase, logger *zap.Logger) *BillingHandler {
	return &BillingHandler{
		paymentUseCase: paymentUseCase,
		logger:         logger,
	}
}

// GetCompanyPayments godoc
// @Summary Get company payments
// @Description Get payments for a company (filtered by user's organization)
// @Tags billing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number"
// @Param page_size query int false "Page size"
// @Param status query string false "Payment status filter"
// @Success 200 {object} PaginatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/company/payments [get]
func (h *BillingHandler) GetCompanyPayments(c *gin.Context) {
	// Get authenticated user info
	_, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	userRoleRaw, exists := c.Get("user_role")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "User role not found",
		})
		return
	}

	var userRole domain.UserRole
	switch v := userRoleRaw.(type) {
	case domain.UserRole:
		userRole = v
	case string:
		userRole = domain.UserRole(v)
	default:
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Invalid user role type",
		})
		return
	}

	if userRole != domain.UserRoleCompany {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error: "Only company users can view company payments",
		})
		return
	}

	orgIDRaw, exists := c.Get("org_id")
	if !exists {
		c.JSON(http.StatusForbidden, ErrorResponse{
			Error: "User must belong to an organization",
		})
		return
	}

	// Parse query parameters
	page := 1
	if p := c.Query("page"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil && parsed > 0 {
			page = parsed
		}
	}

	pageSize := 10
	if ps := c.Query("page_size"); ps != "" {
		if parsed, err := strconv.Atoi(ps); err == nil && parsed > 0 && parsed <= 100 {
			pageSize = parsed
		}
	}

	status := c.Query("status")

	// Parse org_id
	orgUUID, err := uuid.Parse(fmt.Sprint(orgIDRaw))
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid organization ID",
		})
		return
	}

	// Get payments for this company
	payments, total, err := h.paymentUseCase.GetCompanyPayments(orgUUID, page, pageSize, status)
	if err != nil {
		h.logger.Error("Failed to get company payments", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       payments,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}
