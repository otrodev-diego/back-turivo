package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
	"turivo-backend/internal/usecase"
)

type ReservationHandler struct {
	reservationUseCase *usecase.ReservationUseCase
	validator          *validator.Validate
	logger             *zap.Logger
}

func NewReservationHandler(reservationUseCase *usecase.ReservationUseCase, validator *validator.Validate, logger *zap.Logger) *ReservationHandler {
	return &ReservationHandler{
		reservationUseCase: reservationUseCase,
		validator:          validator,
		logger:             logger,
	}
}

// CreateReservationRequest extends the domain request with pricing parameters
type CreateReservationRequest struct {
	domain.CreateReservationRequest
	VehicleType        domain.VehicleType `json:"vehicle_type" validate:"required"`
	HasSpecialLanguage bool               `json:"has_special_language"`
	Stops              int                `json:"stops" validate:"min=0"`
}

// ListReservations godoc
// @Summary List reservations
// @Description Get paginated list of reservations
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string false "Search query"
// @Param status query string false "Filter by status" Enums(ACTIVA,PROGRAMADA,COMPLETADA,CANCELADA)
// @Param from query string false "Filter from date (RFC3339)"
// @Param to query string false "Filter to date (RFC3339)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param sort query string false "Sort by field" Enums(datetime,created_at)
// @Success 200 {object} PaginatedResponse{data=[]domain.Reservation}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/reservations [get]
func (h *ReservationHandler) ListReservations(c *gin.Context) {
	req := domain.ListReservationsRequest{
		Page:     1,
		PageSize: 20,
		Sort:     "datetime",
	}

	if q := c.Query("q"); q != "" {
		req.Query = &q
	}

	if status := c.Query("status"); status != "" {
		reservationStatus := domain.ReservationStatus(status)
		req.Status = &reservationStatus
	}

	if from := c.Query("from"); from != "" {
		if t, err := time.Parse(time.RFC3339, from); err == nil {
			req.From = &t
		}
	}

	if to := c.Query("to"); to != "" {
		if t, err := time.Parse(time.RFC3339, to); err == nil {
			req.To = &t
		}
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

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for list reservations", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	reservations, total, err := h.reservationUseCase.ListReservations(req)
	if err != nil {
		h.logger.Error("Failed to list reservations", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	totalPages := (total + req.PageSize - 1) / req.PageSize

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       reservations,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// CreateReservation godoc
// @Summary Create reservation
// @Description Create a new reservation with automatic price calculation
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body CreateReservationRequest true "Reservation data"
// @Success 201 {object} domain.Reservation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/reservations [post]
func (h *ReservationHandler) CreateReservation(c *gin.Context) {
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

	var req CreateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for create reservation", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	// Set user ID from authenticated user
	req.UserID = &userUUID

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for create reservation", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	reservation, err := h.reservationUseCase.CreateReservation(
		req.CreateReservationRequest,
		req.VehicleType,
		req.HasSpecialLanguage,
		req.Stops,
	)
	if err != nil {
		switch err {
		case domain.ErrReservationPastDate:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Reservation date cannot be in the past",
			})
		default:
			h.logger.Error("Failed to create reservation", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, reservation)
}

// GetReservation godoc
// @Summary Get reservation
// @Description Get reservation by ID
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Reservation ID"
// @Success 200 {object} domain.Reservation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/reservations/{id} [get]
func (h *ReservationHandler) GetReservation(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Reservation ID is required",
		})
		return
	}

	reservation, err := h.reservationUseCase.GetReservationByID(id)
	if err != nil {
		switch err {
		case domain.ErrReservationNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Reservation not found",
			})
		default:
			h.logger.Error("Failed to get reservation", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// UpdateReservation godoc
// @Summary Update reservation
// @Description Update reservation by ID
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Reservation ID"
// @Param request body domain.UpdateReservationRequest true "Reservation data"
// @Success 200 {object} domain.Reservation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/reservations/{id} [patch]
func (h *ReservationHandler) UpdateReservation(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Reservation ID is required",
		})
		return
	}

	var req domain.UpdateReservationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for update reservation", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for update reservation", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	reservation, err := h.reservationUseCase.UpdateReservation(id, req)
	if err != nil {
		switch err {
		case domain.ErrReservationNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Reservation not found",
			})
		case domain.ErrReservationPastDate:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Reservation date cannot be in the past",
			})
		case domain.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Cannot modify completed or cancelled reservation",
			})
		default:
			h.logger.Error("Failed to update reservation", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// ChangeReservationStatus godoc
// @Summary Change reservation status
// @Description Change the status of a reservation
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Reservation ID"
// @Param request body domain.ChangeReservationStatusRequest true "Status change data"
// @Success 200 {object} domain.Reservation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/reservations/{id}/status [patch]
func (h *ReservationHandler) ChangeReservationStatus(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Reservation ID is required",
		})
		return
	}

	var req domain.ChangeReservationStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for change reservation status", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for change reservation status", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	reservation, err := h.reservationUseCase.ChangeReservationStatus(id, req)
	if err != nil {
		switch err {
		case domain.ErrReservationNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Reservation not found",
			})
		case domain.ErrInvalidStatusTransition:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid status transition",
			})
		default:
			h.logger.Error("Failed to change reservation status", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// AssignDriver godoc
// @Summary Assign driver to reservation
// @Description Assign a driver to a specific reservation
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Reservation ID"
// @Param request body domain.AssignDriverRequest true "Driver assignment data"
// @Success 200 {object} domain.Reservation
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/reservations/{id}/driver [patch]
func (h *ReservationHandler) AssignDriver(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Reservation ID is required",
		})
		return
	}

	var req domain.AssignDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for assign driver", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for assign driver", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	reservation, err := h.reservationUseCase.AssignDriver(id, req.DriverID)
	if err != nil {
		switch err {
		case domain.ErrReservationNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Reservation not found",
			})
		case domain.ErrDriverNotFound:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Driver not found",
			})
		case domain.ErrInvalidInput:
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Cannot assign driver to completed or cancelled reservation",
			})
		default:
			h.logger.Error("Failed to assign driver", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, reservation)
}

// GetReservationTimeline godoc
// @Summary Get reservation timeline
// @Description Get timeline events for a reservation
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Reservation ID"
// @Success 200 {object} []domain.TimelineEvent
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/reservations/{id}/timeline [get]
func (h *ReservationHandler) GetReservationTimeline(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Reservation ID is required",
		})
		return
	}

	timeline, err := h.reservationUseCase.GetReservationTimeline(id)
	if err != nil {
		switch err {
		case domain.ErrReservationNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Reservation not found",
			})
		default:
			h.logger.Error("Failed to get reservation timeline", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, timeline)
}

// AddTimelineEvent godoc
// @Summary Add timeline event
// @Description Add a new event to reservation timeline
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Reservation ID"
// @Param request body domain.CreateTimelineEventRequest true "Timeline event data"
// @Success 201 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/reservations/{id}/timeline [post]
func (h *ReservationHandler) AddTimelineEvent(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Reservation ID is required",
		})
		return
	}

	var req domain.CreateTimelineEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for add timeline event", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for add timeline event", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	if err := h.reservationUseCase.AddTimelineEvent(id, req); err != nil {
		switch err {
		case domain.ErrReservationNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Reservation not found",
			})
		default:
			h.logger.Error("Failed to add timeline event", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, MessageResponse{
		Message: "Timeline event added successfully",
	})
}

// GetMyReservations godoc
// @Summary Get current user's reservations
// @Description Get paginated list of reservations for the authenticated user
// @Tags reservations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param status query string false "Filter by status" Enums(ACTIVA,PROGRAMADA,COMPLETADA,CANCELADA)
// @Param from query string false "Filter from date (RFC3339)"
// @Param to query string false "Filter to date (RFC3339)"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param sort query string false "Sort by field" Enums(datetime,created_at)
// @Success 200 {object} PaginatedResponse{data=[]domain.Reservation}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/reservations/my [get]
func (h *ReservationHandler) GetMyReservations(c *gin.Context) {
	// Get user ID from JWT token
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

	// Parse query parameters
	page := 1
	pageSize := 20
	statusStr := c.Query("status")
	sort := c.Query("sort")

	if pageStr := c.Query("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}
	if pageSizeStr := c.Query("page_size"); pageSizeStr != "" {
		if ps, err := strconv.Atoi(pageSizeStr); err == nil && ps > 0 && ps <= 100 {
			pageSize = ps
		}
	}

	// Convert status string to domain type
	var status *domain.ReservationStatus
	if statusStr != "" {
		// Handle multiple statuses separated by comma
		// For simplicity, take the first one for now
		statusParts := strings.Split(statusStr, ",")
		if len(statusParts) > 0 && statusParts[0] != "" {
			reservationStatus := domain.ReservationStatus(statusParts[0])
			status = &reservationStatus
		}
	}

	// Create list request with user filter
	req := domain.ListReservationsRequest{
		Page:     page,
		PageSize: pageSize,
		Status:   status,
		Sort:     sort,
		UserID:   &userUUID, // Filter by current user
	}

	reservations, total, err := h.reservationUseCase.ListReservations(req)
	if err != nil {
		h.logger.Error("Failed to list user reservations", zap.Error(err), zap.String("user_id", userUUID.String()))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to retrieve reservations",
		})
		return
	}

	totalPages := (total + pageSize - 1) / pageSize

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       reservations,
		Page:       page,
		PageSize:   pageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}
