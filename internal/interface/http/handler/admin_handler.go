package handler

import (
	"net/http"
	"time"

	"turivo-backend/internal/domain"
	"turivo-backend/internal/infrastructure/repository"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type AdminHandler struct {
	userRepo        repository.UserRepository
	driverRepo      repository.DriverRepository
	vehicleRepo     repository.VehicleRepository
	reservationRepo repository.ReservationRepository
	companyRepo     repository.CompanyRepository
	logger          *zap.Logger
}

func NewAdminHandler(
	userRepo repository.UserRepository,
	driverRepo repository.DriverRepository,
	vehicleRepo repository.VehicleRepository,
	reservationRepo repository.ReservationRepository,
	companyRepo repository.CompanyRepository,
	logger *zap.Logger,
) *AdminHandler {
	return &AdminHandler{
		userRepo:        userRepo,
		driverRepo:      driverRepo,
		vehicleRepo:     vehicleRepo,
		reservationRepo: reservationRepo,
		companyRepo:     companyRepo,
		logger:          logger,
	}
}

// AdminDashboardResponse represents the admin dashboard data
type AdminDashboardResponse struct {
	KPIs struct {
		TotalUsers            int `json:"total_users"`
		TotalDrivers          int `json:"total_drivers"`
		TotalVehicles         int `json:"total_vehicles"`
		TotalReservations     int `json:"total_reservations"`
		ActiveReservations    int `json:"active_reservations"`
		CompletedReservations int `json:"completed_reservations"`
		TotalCompanies        int `json:"total_companies"`
	} `json:"kpis"`
	RecentReservations []domain.Reservation `json:"recent_reservations"`
	VehicleStats       struct {
		ByType   map[string]int `json:"by_type"`
		ByStatus map[string]int `json:"by_status"`
	} `json:"vehicle_stats"`
	ReservationStats struct {
		ByStatus map[string]int `json:"by_status"`
		ByMonth  []MonthlyStats `json:"by_month"`
	} `json:"reservation_stats"`
}

type MonthlyStats struct {
	Month string `json:"month"`
	Count int    `json:"count"`
}

// GetAdminDashboard godoc
// @Summary Get admin dashboard data
// @Description Get comprehensive dashboard data for admin users
// @Tags admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} AdminDashboardResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/admin/dashboard [get]
func (h *AdminHandler) GetAdminDashboard(c *gin.Context) {
	h.logger.Info("GetAdminDashboard endpoint called")

	// Get total counts
	users, _, err := h.userRepo.List(domain.ListUsersRequest{
		Page:     1,
		PageSize: 1000,
	})
	if err != nil {
		h.logger.Error("Failed to get users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get users"})
		return
	}

	drivers, _, err := h.driverRepo.List(domain.ListDriversRequest{
		Page:     1,
		PageSize: 1000,
	})
	if err != nil {
		h.logger.Error("Failed to get drivers", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get drivers"})
		return
	}

	vehicles, _, err := h.vehicleRepo.List(domain.ListVehiclesRequest{
		Page:     1,
		PageSize: 1000,
	})
	if err != nil {
		h.logger.Error("Failed to get vehicles", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get vehicles"})
		return
	}

	reservations, _, err := h.reservationRepo.List(domain.ListReservationsRequest{
		Page:     1,
		PageSize: 1000,
	})
	if err != nil {
		h.logger.Error("Failed to get reservations", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get reservations"})
		return
	}

	companies, _, err := h.companyRepo.List(domain.ListCompaniesRequest{
		Page:     1,
		PageSize: 1000,
	})
	if err != nil {
		h.logger.Error("Failed to get companies", zap.Error(err))
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to get companies"})
		return
	}

	// Calculate stats
	activeReservations := 0
	completedReservations := 0
	reservationByStatus := make(map[string]int)

	for _, reservation := range reservations {
		reservationByStatus[string(reservation.Status)]++
		if reservation.Status == domain.ReservationStatusProgramada {
			activeReservations++
		} else if reservation.Status == domain.ReservationStatusCompletada {
			completedReservations++
		}
	}

	// Vehicle stats
	vehicleByType := make(map[string]int)
	vehicleByStatus := make(map[string]int)
	for _, vehicle := range vehicles {
		vehicleByType[string(vehicle.Type)]++
		vehicleByStatus[string(vehicle.Status)]++
	}

	// Get recent reservations (last 10)
	recentReservations := make([]domain.Reservation, len(reservations))
	for i, res := range reservations {
		recentReservations[i] = *res
	}
	if len(recentReservations) > 10 {
		recentReservations = recentReservations[:10]
	}

	// Monthly stats (last 6 months)
	monthlyStats := h.calculateMonthlyStats(recentReservations)

	response := AdminDashboardResponse{
		KPIs: struct {
			TotalUsers            int `json:"total_users"`
			TotalDrivers          int `json:"total_drivers"`
			TotalVehicles         int `json:"total_vehicles"`
			TotalReservations     int `json:"total_reservations"`
			ActiveReservations    int `json:"active_reservations"`
			CompletedReservations int `json:"completed_reservations"`
			TotalCompanies        int `json:"total_companies"`
		}{
			TotalUsers:            len(users),
			TotalDrivers:          len(drivers),
			TotalVehicles:         len(vehicles),
			TotalReservations:     len(reservations),
			ActiveReservations:    activeReservations,
			CompletedReservations: completedReservations,
			TotalCompanies:        len(companies),
		},
		RecentReservations: recentReservations,
		VehicleStats: struct {
			ByType   map[string]int `json:"by_type"`
			ByStatus map[string]int `json:"by_status"`
		}{
			ByType:   vehicleByType,
			ByStatus: vehicleByStatus,
		},
		ReservationStats: struct {
			ByStatus map[string]int `json:"by_status"`
			ByMonth  []MonthlyStats `json:"by_month"`
		}{
			ByStatus: reservationByStatus,
			ByMonth:  monthlyStats,
		},
	}

	h.logger.Info("Admin dashboard data retrieved successfully")
	c.JSON(http.StatusOK, response)
}

func (h *AdminHandler) calculateMonthlyStats(reservations []domain.Reservation) []MonthlyStats {
	monthlyCounts := make(map[string]int)

	// Initialize last 6 months
	now := time.Now()
	for i := 5; i >= 0; i-- {
		month := now.AddDate(0, -i, 0)
		monthKey := month.Format("2006-01")
		monthlyCounts[monthKey] = 0
	}

	// Count reservations by month
	for _, reservation := range reservations {
		monthKey := reservation.CreatedAt.Format("2006-01")
		if _, exists := monthlyCounts[monthKey]; exists {
			monthlyCounts[monthKey]++
		}
	}

	// Convert to slice
	var stats []MonthlyStats
	for month, count := range monthlyCounts {
		stats = append(stats, MonthlyStats{
			Month: month,
			Count: count,
		})
	}

	return stats
}
