package http

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"turivo-backend/internal/interface/http/handler"
	"turivo-backend/internal/interface/http/middleware"
	"turivo-backend/internal/usecase"
)

type Router struct {
	authHandler            *handler.AuthHandler
	userHandler            *handler.UserHandler
	driverHandler          *handler.DriverHandler
	reservationHandler     *handler.ReservationHandler
	paymentHandler         *handler.PaymentHandler
	companyHandler         *handler.CompanyHandler
	companyUserHandler     *handler.CompanyHandler // Para operaciones de usuarios de empresa
	vehicleHandler         *handler.VehicleHandler
	supportHandler         *handler.SupportHandler
	authMiddleware         *middleware.AuthMiddleware
	logger                 *zap.Logger
}

func NewRouter(
	authUseCase *usecase.AuthUseCase,
	userUseCase *usecase.UserUseCase,
	driverUseCase *usecase.DriverUseCase,
	reservationUseCase *usecase.ReservationUseCase,
	paymentUseCase *usecase.PaymentUseCase,
	companyUseCase *usecase.CompanyUseCase,
	vehicleUseCase domain.VehicleUseCase,
	authMiddleware *middleware.AuthMiddleware,
	validator *validator.Validate,
	logger *zap.Logger,
) *Router {
	return &Router{
		authHandler:        handler.NewAuthHandler(authUseCase, validator, logger),
		userHandler:        handler.NewUserHandler(userUseCase, validator, logger),
		driverHandler:      handler.NewDriverHandler(driverUseCase, validator, logger),
		reservationHandler: handler.NewReservationHandler(reservationUseCase, validator, logger),
		paymentHandler:     handler.NewPaymentHandler(paymentUseCase, validator, logger),
		companyHandler:     handler.NewCompanyHandler(companyUseCase, validator, logger),
		vehicleHandler:     handler.NewVehicleHandler(vehicleUseCase, validator, logger),
		// supportHandler will be initialized in the main app setup
		authMiddleware: authMiddleware,
		logger:         logger,
	}
}

func (r *Router) SetSupportHandler(supportHandler *handler.SupportHandler) {
	r.supportHandler = supportHandler
}

func (r *Router) SetupRoutes(engine *gin.Engine) {
	// API v1 routes
	v1 := engine.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", r.authHandler.Login)
			auth.POST("/refresh", r.authHandler.RefreshToken)
			auth.POST("/logout", r.authHandler.Logout)
			auth.POST("/complete-registration", r.userHandler.CompleteRegistration)
			auth.GET("/validate-token", r.userHandler.ValidateRegistrationToken)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(r.authMiddleware.RequireAuth())
		{
			// Users routes (Admin can see all, User and Company can see their org)
			users := protected.Group("/users")
			users.Use(r.authMiddleware.RequireRole("ADMIN", "USER", "COMPANY"))
			// users.Use(r.authMiddleware.RequireOrgScope())
			{
				users.GET("", r.userHandler.ListUsers)
				// All authenticated users can get user details
				users.GET("/:id", r.userHandler.GetUser)
			}

			// Admin-only user operations (separate group to avoid middleware conflicts)
			adminUsers := protected.Group("/users")
			adminUsers.Use(r.authMiddleware.RequireRole("ADMIN"))
			{
				adminUsers.POST("", r.userHandler.CreateUser)
				adminUsers.POST("/invite", r.userHandler.CreateUserInvitation)
				adminUsers.PATCH("/:id", r.userHandler.UpdateUser)
				adminUsers.DELETE("/:id", r.userHandler.DeleteUser)
			}

			// Drivers routes (Admin and Company)
			drivers := protected.Group("/drivers")
			drivers.Use(r.authMiddleware.RequireRole("ADMIN", "COMPANY"))
			{
				drivers.GET("", r.driverHandler.ListDrivers)
				drivers.POST("", r.driverHandler.CreateDriver)
				drivers.GET("/:id", r.driverHandler.GetDriver)
				drivers.PATCH("/:id", r.driverHandler.UpdateDriver)
				drivers.DELETE("/:id", r.driverHandler.DeleteDriver)
				drivers.GET("/:id/kpis", r.driverHandler.GetDriverKPIs)
			}

			// Reservations routes (All authenticated users)
			reservations := protected.Group("/reservations")
			reservations.Use(r.authMiddleware.RequireOrgScope())
			{
				reservations.GET("", r.reservationHandler.ListReservations)
				reservations.GET("/my", r.reservationHandler.GetMyReservations) // User's own reservations
				reservations.POST("", r.reservationHandler.CreateReservation)
				// Specific routes MUST come before generic /:id routes
				reservations.PATCH("/:id/status", r.reservationHandler.ChangeReservationStatus)
				reservations.PATCH("/:id/driver", r.reservationHandler.AssignDriver)
				reservations.GET("/:id/test", func(c *gin.Context) { c.JSON(200, gin.H{"test": "router.go", "id": c.Param("id")}) })
				reservations.GET("/:id/timeline", r.reservationHandler.GetReservationTimeline)
				reservations.POST("/:id/timeline", r.reservationHandler.AddTimelineEvent)
				// Generic routes come last
				reservations.GET("/:id", r.reservationHandler.GetReservation)
				reservations.PATCH("/:id", r.reservationHandler.UpdateReservation)
			}

			// Payments routes (All authenticated users)
			payments := protected.Group("/payments")
			{
				payments.POST("", r.paymentHandler.CreatePayment)
				payments.GET("/:id", r.paymentHandler.GetPayment)
				payments.POST("/:id/simulate", r.paymentHandler.SimulatePayment)
			}

			// Companies routes (Admin can see all, User and Company can see their org)
			companies := protected.Group("/companies")
			companies.Use(r.authMiddleware.RequireRole("ADMIN", "USER", "COMPANY"))
			// companies.Use(r.authMiddleware.RequireOrgScope())
			{
				companies.GET("", r.companyHandler.ListCompanies)
				companies.GET("/:id", r.companyHandler.GetCompany)
			}

			// Admin-only company operations (separate group to avoid middleware conflicts)
			adminCompanies := protected.Group("/companies")
			adminCompanies.Use(r.authMiddleware.RequireRole("ADMIN"))
			{
				adminCompanies.POST("", r.companyHandler.CreateCompany)
				adminCompanies.PUT("/:id", r.companyHandler.UpdateCompany)
				adminCompanies.DELETE("/:id", r.companyHandler.DeleteCompany)
			}

			// Vehicles routes (Admin and Company)
			vehicles := protected.Group("/vehicles")
			vehicles.Use(r.authMiddleware.RequireRole("ADMIN", "COMPANY"))
			{
				vehicles.GET("", r.vehicleHandler.ListVehicles)
				vehicles.GET("/:id", r.vehicleHandler.GetVehicle)
				vehicles.POST("", r.vehicleHandler.CreateVehicle)
				vehicles.PUT("/:id", r.vehicleHandler.UpdateVehicle)
				vehicles.DELETE("/:id", r.vehicleHandler.DeleteVehicle)
				vehicles.POST("/:id/assign", r.vehicleHandler.AssignVehicleToDriver)
			}

			// Support routes (All authenticated users)
			if r.supportHandler != nil {
				support := protected.Group("/support")
				{
					support.POST("/contact", r.supportHandler.ContactSupport)
				}
			}
		}
	}

	r.logger.Info("Routes configured successfully")
}
