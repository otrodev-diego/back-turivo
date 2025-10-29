package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
	"turivo-backend/internal/usecase"
)

type DriverHandler struct {
	driverUseCase *usecase.DriverUseCase
	validator     *validator.Validate
	logger        *zap.Logger
}

func NewDriverHandler(driverUseCase *usecase.DriverUseCase, validator *validator.Validate, logger *zap.Logger) *DriverHandler {
	return &DriverHandler{
		driverUseCase: driverUseCase,
		validator:     validator,
		logger:        logger,
	}
}

// ListDrivers godoc
// @Summary List drivers
// @Description Get paginated list of drivers
// @Tags drivers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string false "Search query"
// @Param status query string false "Filter by status" Enums(ACTIVE,INACTIVE)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param sort query string false "Sort by field" Enums(name,id,created_at)
// @Success 200 {object} PaginatedResponse{data=[]domain.Driver}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/drivers [get]
func (h *DriverHandler) ListDrivers(c *gin.Context) {
	req := domain.ListDriversRequest{
		Page:     1,
		PageSize: 20,
		Sort:     "created_at",
	}

	if q := c.Query("q"); q != "" {
		req.Query = &q
	}

	if status := c.Query("status"); status != "" {
		driverStatus := domain.DriverStatus(status)
		req.Status = &driverStatus
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
		h.logger.Warn("Validation failed for list drivers", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	drivers, total, err := h.driverUseCase.ListDrivers(req)
	if err != nil {
		h.logger.Error("Failed to list drivers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	totalPages := (total + req.PageSize - 1) / req.PageSize

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       drivers,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// CreateDriver godoc
// @Summary Create driver
// @Description Create a new driver
// @Tags drivers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body domain.CreateDriverRequest true "Driver data"
// @Success 201 {object} domain.Driver
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/drivers [post]
func (h *DriverHandler) CreateDriver(c *gin.Context) {
	var req domain.CreateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for create driver", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for create driver", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	driver, err := h.driverUseCase.CreateDriver(req)
	if err != nil {
		switch err {
		case domain.ErrDriverAlreadyExists:
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Driver already exists",
			})
		default:
			h.logger.Error("Failed to create driver", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusCreated, driver)
}

// GetDriver godoc
// @Summary Get driver
// @Description Get driver by ID
// @Tags drivers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Driver ID"
// @Success 200 {object} domain.Driver
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/drivers/{id} [get]
func (h *DriverHandler) GetDriver(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Driver ID is required",
		})
		return
	}

	driver, err := h.driverUseCase.GetDriverByID(id)
	if err != nil {
		switch err {
		case domain.ErrDriverNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Driver not found",
			})
		default:
			h.logger.Error("Failed to get driver", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, driver)
}

// UpdateDriver godoc
// @Summary Update driver
// @Description Update driver by ID
// @Tags drivers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Driver ID"
// @Param request body domain.UpdateDriverRequest true "Driver data"
// @Success 200 {object} domain.Driver
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/drivers/{id} [patch]
func (h *DriverHandler) UpdateDriver(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Driver ID is required",
		})
		return
	}

	var req domain.UpdateDriverRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid request body for update driver", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request body",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for update driver", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	driver, err := h.driverUseCase.UpdateDriver(id, req)
	if err != nil {
		switch err {
		case domain.ErrDriverNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Driver not found",
			})
		default:
			h.logger.Error("Failed to update driver", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, driver)
}

// DeleteDriver godoc
// @Summary Delete driver
// @Description Delete driver by ID
// @Tags drivers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Driver ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/drivers/{id} [delete]
func (h *DriverHandler) DeleteDriver(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Driver ID is required",
		})
		return
	}

	if err := h.driverUseCase.DeleteDriver(id); err != nil {
		switch err {
		case domain.ErrDriverNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Driver not found",
			})
		default:
			h.logger.Error("Failed to delete driver", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, MessageResponse{
		Message: "Driver deleted successfully",
	})
}

// GetDriverKPIs godoc
// @Summary Get driver KPIs
// @Description Get driver key performance indicators
// @Tags drivers
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Driver ID"
// @Success 200 {object} domain.DriverKPIs
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/drivers/{id}/kpis [get]
func (h *DriverHandler) GetDriverKPIs(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Driver ID is required",
		})
		return
	}

	kpis, err := h.driverUseCase.GetDriverKPIs(id)
	if err != nil {
		switch err {
		case domain.ErrDriverNotFound:
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Driver not found",
			})
		default:
			h.logger.Error("Failed to get driver KPIs", zap.Error(err))
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error: "Internal server error",
			})
		}
		return
	}

	c.JSON(http.StatusOK, kpis)
}
