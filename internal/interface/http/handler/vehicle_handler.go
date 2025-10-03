package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type VehicleHandler struct {
	vehicleUseCase domain.VehicleUseCase
	validator      *validator.Validate
	logger         *zap.Logger
}

func NewVehicleHandler(vehicleUseCase domain.VehicleUseCase, validator *validator.Validate, logger *zap.Logger) *VehicleHandler {
	return &VehicleHandler{
		vehicleUseCase: vehicleUseCase,
		validator:      validator,
		logger:         logger,
	}
}

// ListVehicles godoc
// @Summary List vehicles
// @Description Get paginated list of vehicles
// @Tags vehicles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string false "Search query"
// @Param type query string false "Filter by type" Enums(BUS,VAN,SEDAN,SUV)
// @Param status query string false "Filter by status" Enums(AVAILABLE,ASSIGNED,MAINTENANCE,INACTIVE)
// @Param driver_id query string false "Filter by driver ID"
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param sort query string false "Sort by field" Enums(brand,model,year,created_at)
// @Success 200 {object} PaginatedResponse{data=[]domain.Vehicle}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/vehicles [get]
func (h *VehicleHandler) ListVehicles(c *gin.Context) {
	h.logger.Info("ListVehicles endpoint called")

	// Parse query parameters
	req := domain.ListVehiclesRequest{
		Page:     1,
		PageSize: 20,
		Sort:     "created_at",
	}

	if q := c.Query("q"); q != "" {
		req.Query = &q
	}

	if vehicleType := c.Query("type"); vehicleType != "" {
		vType := domain.VehicleType(vehicleType)
		req.Type = &vType
	}

	if status := c.Query("status"); status != "" {
		vStatus := domain.VehicleStatus(status)
		req.Status = &vStatus
	}

	if driverID := c.Query("driver_id"); driverID != "" {
		req.DriverID = &driverID
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
		h.logger.Warn("Validation failed for list vehicles", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request parameters",
		})
		return
	}

	vehicles, total, err := h.vehicleUseCase.ListVehicles(req)
	if err != nil {
		h.logger.Error("Failed to list vehicles", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to list vehicles",
		})
		return
	}

	totalPages := (total + req.PageSize - 1) / req.PageSize

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       vehicles,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// GetVehicle godoc
// @Summary Get vehicle by ID
// @Description Get a vehicle by its ID
// @Tags vehicles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Vehicle ID"
// @Success 200 {object} domain.Vehicle
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/vehicles/{id} [get]
func (h *VehicleHandler) GetVehicle(c *gin.Context) {
	id := c.Param("id")

	h.logger.Info("GetVehicle endpoint called", zap.String("id", id))

	vehicle, err := h.vehicleUseCase.GetVehicle(id)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Vehicle not found",
			})
			return
		}
		h.logger.Error("Failed to get vehicle", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to get vehicle",
		})
		return
	}

	c.JSON(http.StatusOK, vehicle)
}

// CreateVehicle godoc
// @Summary Create a new vehicle
// @Description Create a new vehicle
// @Tags vehicles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param vehicle body domain.CreateVehicleRequest true "Vehicle data"
// @Success 201 {object} domain.Vehicle
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/vehicles [post]
func (h *VehicleHandler) CreateVehicle(c *gin.Context) {
	var req domain.CreateVehicleRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Invalid request body - JSON parsing failed", 
			zap.Error(err),
			zap.Any("raw_body", c.Request.Body))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body: " + err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Validation failed: " + err.Error(),
		})
		return
	}

	vehicle, err := h.vehicleUseCase.CreateVehicle(req)
	if err != nil {
		if err == domain.ErrInvalidInput {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid vehicle data",
			})
			return
		}
		h.logger.Error("Failed to create vehicle", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to create vehicle",
		})
		return
	}

	h.logger.Info("Vehicle created successfully", zap.String("id", vehicle.ID))

	c.JSON(http.StatusCreated, vehicle)
}

// UpdateVehicle godoc
// @Summary Update a vehicle
// @Description Update vehicle information
// @Tags vehicles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Vehicle ID"
// @Param vehicle body domain.UpdateVehicleRequest true "Vehicle data to update"
// @Success 200 {object} domain.Vehicle
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/vehicles/{id} [put]
func (h *VehicleHandler) UpdateVehicle(c *gin.Context) {
	id := c.Param("id")

	var req domain.UpdateVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Validation failed: " + err.Error(),
		})
		return
	}

	vehicle, err := h.vehicleUseCase.UpdateVehicle(id, req)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Vehicle not found",
			})
			return
		}
		if err == domain.ErrInvalidInput {
			c.JSON(http.StatusBadRequest, ErrorResponse{
				Error: "Invalid vehicle data",
			})
			return
		}
		h.logger.Error("Failed to update vehicle", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to update vehicle",
		})
		return
	}

	h.logger.Info("Vehicle updated successfully", zap.String("id", id))

	c.JSON(http.StatusOK, vehicle)
}

// DeleteVehicle godoc
// @Summary Delete a vehicle
// @Description Delete a vehicle by ID
// @Tags vehicles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Vehicle ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/vehicles/{id} [delete]
func (h *VehicleHandler) DeleteVehicle(c *gin.Context) {
	id := c.Param("id")

	h.logger.Info("DeleteVehicle endpoint called", zap.String("id", id))

	err := h.vehicleUseCase.DeleteVehicle(id)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Vehicle not found",
			})
			return
		}
		h.logger.Error("Failed to delete vehicle", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: err.Error(),
		})
		return
	}

	h.logger.Info("Vehicle deleted successfully", zap.String("id", id))

	c.JSON(http.StatusNoContent, nil)
}

// AssignVehicleToDriver godoc
// @Summary Assign vehicle to driver
// @Description Assign a vehicle to a driver (or unassign if driver_id is null)
// @Tags vehicles
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Vehicle ID"
// @Param assignment body domain.AssignVehicleRequest true "Assignment data"
// @Success 200 {object} domain.Vehicle
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/vehicles/{id}/assign [post]
func (h *VehicleHandler) AssignVehicleToDriver(c *gin.Context) {
	id := c.Param("id")

	var req domain.AssignVehicleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid request body",
		})
		return
	}

	// Override vehicle ID from path
	req.VehicleID = id

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Validation failed: " + err.Error(),
		})
		return
	}

	vehicle, err := h.vehicleUseCase.AssignVehicleToDriver(req.VehicleID, req.DriverID)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Vehicle or driver not found",
			})
			return
		}
		h.logger.Error("Failed to assign vehicle", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Failed to assign vehicle: " + err.Error(),
		})
		return
	}

	if req.DriverID != nil {
		h.logger.Info("Vehicle assigned successfully",
			zap.String("vehicle_id", id),
			zap.String("driver_id", *req.DriverID))
	} else {
		h.logger.Info("Vehicle unassigned successfully",
			zap.String("vehicle_id", id))
	}

	c.JSON(http.StatusOK, vehicle)
}

