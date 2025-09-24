package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
	"turivo-backend/internal/usecase"
)

type PaymentHandler struct {
	paymentUseCase *usecase.PaymentUseCase
	validator      *validator.Validate
	logger         *zap.Logger
}

func NewPaymentHandler(paymentUseCase *usecase.PaymentUseCase, validator *validator.Validate, logger *zap.Logger) *PaymentHandler {
	return &PaymentHandler{
		paymentUseCase: paymentUseCase,
		validator:      validator,
		logger:         logger,
	}
}

// CreatePayment godoc
// @Summary Create payment
// @Description Create a new payment for a reservation
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreatePaymentRequest true "Payment data"
// @Success 201 {object} domain.Payment
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments [post]
func (h *PaymentHandler) CreatePayment(c *gin.Context) {
	var req domain.CreatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for create payment", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for create payment", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	payment, err := h.paymentUseCase.CreatePayment(req)
	if err != nil {
		switch err {
		case domain.ErrReservationNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Reservation not found",
			})
		case domain.ErrPaymentAlreadyPaid:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Reservation already has approved payment",
			})
		default:
			h.logger.Error("Failed to create payment", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, payment)
}

// GetPayment godoc
// @Summary Get payment
// @Description Get payment by ID
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payment ID"
// @Success 200 {object} domain.Payment
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/{id} [get]
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid payment ID",
		})
		return
	}

	payment, err := h.paymentUseCase.GetPaymentByID(id)
	if err != nil {
		switch err {
		case domain.ErrPaymentNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Payment not found",
			})
		default:
			h.logger.Error("Failed to get payment", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, payment)
}

// SimulatePayment godoc
// @Summary Simulate payment result
// @Description Simulate payment approval or rejection (for testing)
// @Tags payments
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Payment ID"
// @Param request body domain.SimulatePaymentRequest true "Simulation data"
// @Success 200 {object} domain.Payment
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/payments/{id}/simulate [post]
func (h *PaymentHandler) SimulatePayment(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid payment ID",
		})
		return
	}

	var req domain.SimulatePaymentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for simulate payment", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for simulate payment", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	payment, err := h.paymentUseCase.SimulatePayment(id, req.Result)
	if err != nil {
		switch err {
		case domain.ErrPaymentNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Payment not found",
			})
		case domain.ErrPaymentAlreadyPaid:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Payment already processed",
			})
		default:
			h.logger.Error("Failed to simulate payment", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, payment)
}

