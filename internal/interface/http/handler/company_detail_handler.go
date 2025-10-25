package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"turivo-backend/internal/domain"
	"turivo-backend/internal/infrastructure/repository"
)

type CompanyDetailHandler struct {
	companyRepo repository.CompanyRepository
	userRepo    repository.UserRepository
	logger      *zap.Logger
}

func NewCompanyDetailHandler(
	companyRepo repository.CompanyRepository,
	userRepo repository.UserRepository,
	logger *zap.Logger,
) *CompanyDetailHandler {
	return &CompanyDetailHandler{
		companyRepo: companyRepo,
		userRepo:    userRepo,
		logger:      logger,
	}
}

// CompanyDetailResponse represents the company detail with users
type CompanyDetailResponse struct {
	Company domain.Company `json:"company"`
	Users   []domain.User  `json:"users"`
	Stats   struct {
		TotalUsers    int `json:"total_users"`
		ActiveUsers   int `json:"active_users"`
		InactiveUsers int `json:"inactive_users"`
	} `json:"stats"`
}

// GetCompanyDetail godoc
// @Summary Get company detail with users
// @Description Get detailed information about a company including its users
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Company ID"
// @Success 200 {object} CompanyDetailResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/companies/{id}/detail [get]
func (h *CompanyDetailHandler) GetCompanyDetail(c *gin.Context) {
	companyID := c.Param("id")
	h.logger.Info("GetCompanyDetail endpoint called", zap.String("company_id", companyID))

	// Parse company ID to UUID
	companyUUID, err := uuid.Parse(companyID)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "Invalid company ID"})
		return
	}

	// Get company details
	company, err := h.companyRepo.GetByID(companyUUID)
	if err != nil {
		if err == domain.ErrNotFound {
			c.JSON(http.StatusNotFound, ErrorResponse{Error: "Company not found"})
			return
		}
		h.logger.Error("Failed to get company", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get company"})
		return
	}

	// Get users for this company
	users, _, err := h.userRepo.List(domain.ListUsersRequest{
		Page:     1,
		PageSize: 1000,
		OrgID:    &company.ID,
	})
	if err != nil {
		h.logger.Error("Failed to get company users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get company users"})
		return
	}

	// Calculate stats
	activeUsers := 0
	inactiveUsers := 0
	for _, user := range users {
		if user.Status == domain.UserStatusActive {
			activeUsers++
		} else {
			inactiveUsers++
		}
	}

	// Convert []*domain.User to []domain.User
	usersList := make([]domain.User, len(users))
	for i, user := range users {
		usersList[i] = *user
	}

	response := CompanyDetailResponse{
		Company: *company,
		Users:   usersList,
		Stats: struct {
			TotalUsers    int `json:"total_users"`
			ActiveUsers   int `json:"active_users"`
			InactiveUsers int `json:"inactive_users"`
		}{
			TotalUsers:    len(users),
			ActiveUsers:   activeUsers,
			InactiveUsers: inactiveUsers,
		},
	}

	h.logger.Info("Company detail retrieved successfully", zap.String("company_id", companyID))
	c.JSON(http.StatusOK, response)
}
