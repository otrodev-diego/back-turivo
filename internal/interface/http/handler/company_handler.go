package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
	"turivo-backend/internal/usecase"
)

type CompanyHandler struct {
	companyUseCase *usecase.CompanyUseCase
	validator      *validator.Validate
	logger         *zap.Logger
}

func NewCompanyHandler(companyUseCase *usecase.CompanyUseCase, validator *validator.Validate, logger *zap.Logger) *CompanyHandler {
	return &CompanyHandler{
		companyUseCase: companyUseCase,
		validator:      validator,
		logger:         logger,
	}
}

// ListCompanies godoc
// @Summary List companies
// @Description Get paginated list of companies
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param q query string false "Search query"
// @Param status query string false "Filter by status" Enums(ACTIVE,SUSPENDED)
// @Param sector query string false "Filter by sector" Enums(HOTEL,MINERIA,TURISMO)
// @Param page query int false "Page number" default(1)
// @Param page_size query int false "Page size" default(20)
// @Param sort query string false "Sort by field" Enums(name,rut,created_at)
// @Success 200 {object} PaginatedResponse{data=[]domain.Company}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/companies [get]
func (h *CompanyHandler) ListCompanies(c *gin.Context) {
	h.logger.Info("ListCompanies endpoint called")

	// Parse query parameters
	req := domain.ListCompaniesRequest{
		Page:     1,
		PageSize: 20,
		Sort:     "created_at",
	}

	// Parse query parameters
	if q := c.Query("q"); q != "" {
		req.Query = &q
	}

	if status := c.Query("status"); status != "" {
		companyStatus := domain.CompanyStatus(status)
		req.Status = &companyStatus
	}

	if sector := c.Query("sector"); sector != "" {
		companySector := domain.CompanySector(sector)
		req.Sector = &companySector
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
		h.logger.Info("Filtering companies by organization", zap.String("org_id", req.OrgID.String()))
	} else {
		h.logger.Info("No organization scope - showing all companies (admin)")
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for list companies", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	companies, total, err := h.companyUseCase.ListCompanies(req)
	if err != nil {
		h.logger.Error("Failed to list companies", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	totalPages := (total + req.PageSize - 1) / req.PageSize

	c.JSON(http.StatusOK, PaginatedResponse{
		Data:       companies,
		Page:       req.Page,
		PageSize:   req.PageSize,
		Total:      total,
		TotalPages: totalPages,
	})
}

// CreateCompany godoc
// @Summary Create a new company
// @Description Create a new company with the provided information
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param company body domain.CreateCompanyRequest true "Company data"
// @Success 201 {object} domain.Company
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/companies [post]
func (h *CompanyHandler) CreateCompany(c *gin.Context) {
	var req domain.CreateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid JSON in create company request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid JSON",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for create company", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	company, err := h.companyUseCase.CreateCompany(req)
	if err != nil {
		h.logger.Error("Failed to create company", zap.Error(err))
		if err == domain.ErrAlreadyExists {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Company with this RUT already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusCreated, company)
}

// GetCompany godoc
// @Summary Get company by ID
// @Description Get a specific company by its ID
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Success 200 {object} domain.Company
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/companies/{id} [get]
func (h *CompanyHandler) GetCompany(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Warn("Invalid company ID format", zap.String("id", idParam))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid company ID format",
		})
		return
	}

	company, err := h.companyUseCase.GetCompanyByID(id)
	if err != nil {
		h.logger.Error("Failed to get company", zap.Error(err), zap.String("id", id.String()))
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Company not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, company)
}

// UpdateCompany godoc
// @Summary Update company
// @Description Update a company's information
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Param company body domain.UpdateCompanyRequest true "Company update data"
// @Success 200 {object} domain.Company
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/companies/{id} [put]
func (h *CompanyHandler) UpdateCompany(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Warn("Invalid company ID format", zap.String("id", idParam))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid company ID format",
		})
		return
	}

	var req domain.UpdateCompanyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Warn("Invalid JSON in update company request", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid JSON",
			Details: err.Error(),
		})
		return
	}

	if err := h.validator.Struct(req); err != nil {
		h.logger.Warn("Validation failed for update company", zap.Error(err))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	company, err := h.companyUseCase.UpdateCompany(id, req)
	if err != nil {
		h.logger.Error("Failed to update company", zap.Error(err), zap.String("id", id.String()))
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Company not found",
			})
			return
		}
		if err == domain.ErrAlreadyExists {
			c.JSON(http.StatusConflict, ErrorResponse{
				Error: "Company with this RUT already exists",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	c.JSON(http.StatusOK, company)
}

// DeleteCompany godoc
// @Summary Delete company
// @Description Delete a company by ID
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/companies/{id} [delete]
func (h *CompanyHandler) DeleteCompany(c *gin.Context) {
	idParam := c.Param("id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		h.logger.Warn("Invalid company ID format", zap.String("id", idParam))
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error: "Invalid company ID format",
		})
		return
	}

	err = h.companyUseCase.DeleteCompany(id)
	if err != nil {
		h.logger.Error("Failed to delete company", zap.Error(err), zap.String("id", id.String()))
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error: "Company not found",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error: "Internal server error",
		})
		return
	}

	c.Status(http.StatusNoContent)
}
