package handler

import (
	"net/http"

	"turivo-backend/internal/domain"
	"turivo-backend/internal/usecase"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type DriverDashboardHandler struct {
	driverUseCase *usecase.DriverUseCase
	logger        *zap.Logger
}

func NewDriverDashboardHandler(driverUseCase *usecase.DriverUseCase, logger *zap.Logger) *DriverDashboardHandler {
	return &DriverDashboardHandler{
		driverUseCase: driverUseCase,
		logger:        logger,
	}
}

// GetDriverStats godoc
// @Summary Get driver statistics
// @Description Get statistics for the authenticated driver
// @Tags Driver
// @Accept json
// @Produce json
// @Success 200 {object} DriverStatsResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/driver/stats [get]
func (h *DriverDashboardHandler) GetDriverStats(c *gin.Context) {
	// Get authenticated user info
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	userID := userIDRaw.(uuid.UUID).String()
	h.logger.Info("Getting driver stats", zap.String("user_id", userID))

	// Get driver by user ID
	driver, err := h.driverUseCase.GetDriverByUserID(userID)
	if err != nil {
		h.logger.Error("Failed to get driver by user ID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get driver information",
		})
		return
	}

	// Get driver KPIs (handle case when no KPIs exist)
	var stats DriverStatsResponse
	kpis, err := h.driverUseCase.GetDriverKPIs(driver.ID)
	if err != nil {
		h.logger.Warn("No KPIs found for driver, using default values", zap.Error(err))
		// Use default values when no KPIs exist
		stats = DriverStatsResponse{
			TotalTrips:     0,
			CompletedTrips: 0,
			PendingTrips:   0,
			TotalEarnings:  0,
			AverageRating:  0,
			TotalDistance:  0,
		}
	} else {
		// Calculate additional stats from KPIs
		stats = DriverStatsResponse{
			TotalTrips:     kpis.TotalTrips,
			CompletedTrips: kpis.TotalTrips,     // Assuming all trips are completed for now
			PendingTrips:   0,                   // Will be calculated from actual trips
			TotalEarnings:  kpis.TotalKM * 1200, // Example calculation: 1200 CLP per km
			AverageRating:  kpis.AverageRating,
			TotalDistance:  kpis.TotalKM,
		}
	}

	h.logger.Info("Driver stats retrieved successfully", zap.String("driver_id", driver.ID))
	c.JSON(http.StatusOK, stats)
}

// GetDriverTrips godoc
// @Summary Get driver trips
// @Description Get trips assigned to the authenticated driver
// @Tags Driver
// @Accept json
// @Produce json
// @Success 200 {object} []TripResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/driver/trips [get]
func (h *DriverDashboardHandler) GetDriverTrips(c *gin.Context) {
	// Get authenticated user info
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		h.logger.Warn("❌ User ID not found in context")
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	userID := userIDRaw.(uuid.UUID).String()
	h.logger.Info("🚗 Getting driver trips", zap.String("user_id", userID))

	// Get driver by user ID
	driver, err := h.driverUseCase.GetDriverByUserID(userID)
	if err != nil {
		h.logger.Error("❌ Failed to get driver by user ID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get driver information",
		})
		return
	}
	h.logger.Info("✅ Driver found", zap.String("driver_id", driver.ID))

	// Get driver trips
	trips, err := h.driverUseCase.GetDriverTrips(driver.ID)
	if err != nil {
		h.logger.Error("❌ Failed to get driver trips", zap.Error(err), zap.String("driver_id", driver.ID))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get driver trips",
		})
		return
	}

	h.logger.Info("✅ Driver trips retrieved successfully", zap.String("driver_id", driver.ID), zap.Int("count", len(trips)))
	c.JSON(http.StatusOK, trips)
}

// GetDriverVehicle godoc
// @Summary Get driver vehicle
// @Description Get vehicle assigned to the authenticated driver
// @Tags Driver
// @Accept json
// @Produce json
// @Success 200 {object} VehicleResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/driver/vehicle [get]
func (h *DriverDashboardHandler) GetDriverVehicle(c *gin.Context) {
	// Get authenticated user info
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	userID := userIDRaw.(uuid.UUID).String()
	h.logger.Info("Getting driver vehicle", zap.String("user_id", userID))

	// Get driver by user ID
	driver, err := h.driverUseCase.GetDriverByUserID(userID)
	if err != nil {
		h.logger.Error("Failed to get driver by user ID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get driver information",
		})
		return
	}

	// Get driver vehicle
	vehicle, err := h.driverUseCase.GetDriverVehicle(driver.ID)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "No vehicle assigned to this driver",
			})
			return
		}
		h.logger.Error("Failed to get driver vehicle", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get driver vehicle",
		})
		return
	}

	h.logger.Info("Driver vehicle retrieved successfully", zap.String("driver_id", driver.ID))
	c.JSON(http.StatusOK, vehicle)
}

// GetDriverProfile godoc
// @Summary Get driver profile
// @Description Get profile information for the authenticated driver
// @Tags Driver
// @Accept json
// @Produce json
// @Success 200 {object} DriverProfileResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/driver/profile [get]
func (h *DriverDashboardHandler) GetDriverProfile(c *gin.Context) {
	// Get authenticated user info
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	userID := userIDRaw.(uuid.UUID).String()
	h.logger.Info("Getting driver profile", zap.String("user_id", userID))

	// Get driver by user ID
	driver, err := h.driverUseCase.GetDriverByUserID(userID)
	if err != nil {
		h.logger.Error("Failed to get driver by user ID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get driver information",
		})
		return
	}

	// Get driver KPIs for stats
	kpis, err := h.driverUseCase.GetDriverKPIs(driver.ID)
	if err != nil {
		h.logger.Warn("Failed to get driver KPIs for profile", zap.Error(err))
		// Continue without KPIs
	}

	// Convert birth date to string if present
	var birthDateStr *string
	if driver.BirthDate != nil {
		birthDate := driver.BirthDate.Format("2006-01-02")
		birthDateStr = &birthDate
	}

	profile := DriverProfileResponse{
		ID:        driver.ID,
		FirstName: driver.FirstName,
		LastName:  driver.LastName,
		Email:     driver.Email,
		Phone:     driver.Phone,
		RutOrDNI:  driver.RutOrDNI,
		BirthDate: birthDateStr,
		PhotoURL:  driver.PhotoURL,
		Status:    string(driver.Status),
		Vehicle:   driver.Vehicle,
		Company:   driver.Company,
		Stats: &DriverStats{
			TotalTrips:    kpis.TotalTrips,
			AverageRating: kpis.AverageRating,
			TotalEarnings: kpis.TotalKM * 1200, // Example calculation
		},
	}

	h.logger.Info("Driver profile retrieved successfully", zap.String("driver_id", driver.ID))
	c.JSON(http.StatusOK, profile)
}

// UpdateTripStatus godoc
// @Summary Update trip status
// @Description Update the status of a trip for the authenticated driver
// @Tags Driver
// @Accept json
// @Produce json
// @Param tripId path string true "Trip ID"
// @Param request body UpdateTripStatusRequest true "Update trip status request"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /api/v1/driver/trips/{tripId}/status [patch]
func (h *DriverDashboardHandler) UpdateTripStatus(c *gin.Context) {
	tripID := c.Param("tripId")

	var req UpdateTripStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Get authenticated user info
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	userID := userIDRaw.(uuid.UUID).String()
	h.logger.Info("Updating trip status", zap.String("user_id", userID), zap.String("trip_id", tripID), zap.String("status", req.Status))

	// Get driver by user ID
	driver, err := h.driverUseCase.GetDriverByUserID(userID)
	if err != nil {
		h.logger.Error("Failed to get driver by user ID", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get driver information",
		})
		return
	}

	// Update trip status
	err = h.driverUseCase.UpdateTripStatus(driver.ID, tripID, req.Status)
	if err != nil {
		h.logger.Error("Failed to update trip status", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to update trip status",
		})
		return
	}

	h.logger.Info("Trip status updated successfully", zap.String("driver_id", driver.ID), zap.String("trip_id", tripID))
	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Trip status updated successfully",
	})
}

// Response types
type DriverStatsResponse struct {
	TotalTrips     int     `json:"total_trips"`
	CompletedTrips int     `json:"completed_trips"`
	PendingTrips   int     `json:"pending_trips"`
	TotalEarnings  float64 `json:"total_earnings"`
	AverageRating  float64 `json:"average_rating"`
	TotalDistance  float64 `json:"total_distance"`
}

type DriverProfileResponse struct {
	ID        string          `json:"id"`
	FirstName string          `json:"first_name"`
	LastName  string          `json:"last_name"`
	Email     *string         `json:"email,omitempty"`
	Phone     *string         `json:"phone,omitempty"`
	RutOrDNI  string          `json:"rut_or_dni"`
	BirthDate *string         `json:"birth_date,omitempty"`
	PhotoURL  *string         `json:"photo_url,omitempty"`
	Status    string          `json:"status"`
	Vehicle   *domain.Vehicle `json:"vehicle,omitempty"`
	Company   *domain.Company `json:"company,omitempty"`
	Stats     *DriverStats    `json:"stats,omitempty"`
}

type DriverStats struct {
	TotalTrips    int     `json:"total_trips"`
	AverageRating float64 `json:"average_rating"`
	TotalEarnings float64 `json:"total_earnings"`
}

type UpdateTripStatusRequest struct {
	Status string `json:"status" validate:"required,oneof=PENDING IN_PROGRESS COMPLETED CANCELLED"`
}

type SuccessResponse struct {
	Message string `json:"message"`
}
