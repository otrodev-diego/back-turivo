package router

import (
	"github.com/gin-gonic/gin"

	"turivo-backend/internal/interface/http/handler"
	"turivo-backend/internal/interface/http/middleware"
)

type RouteHandlers struct {
	Auth        *handler.AuthHandler
	User        *handler.UserHandler
	Driver      *handler.DriverHandler
	Reservation *handler.ReservationHandler
	Payment     *handler.PaymentHandler
	Company     *handler.CompanyHandler
	Vehicle     *handler.VehicleHandler
	Support     *handler.SupportHandler
}

func SetupRoutes(engine *gin.Engine, handlers RouteHandlers, authMiddleware *middleware.AuthMiddleware) {
	// API v1 routes
	v1 := engine.Group("/api/v1")
	{
		// Public routes (no authentication required)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", handlers.Auth.Login)
			auth.POST("/refresh", handlers.Auth.RefreshToken)
			auth.POST("/logout", handlers.Auth.Logout)
			auth.POST("/complete-registration", handlers.User.CompleteRegistration)
			auth.GET("/validate-token", handlers.User.ValidateRegistrationToken)
		}

		// Protected routes (authentication required)
		protected := v1.Group("")
		protected.Use(authMiddleware.RequireAuth())
		{
			// Users routes (Admin can see all, User and Company can see their org)
			users := protected.Group("/users")
			users.Use(authMiddleware.RequireRole("ADMIN", "USER", "COMPANY"))
			users.Use(authMiddleware.RequireOrgScope())
			{
				users.GET("", handlers.User.ListUsers)
				users.GET("/:id", handlers.User.GetUser)
			}

			// Admin-only user operations
			adminUsers := protected.Group("/users")
			adminUsers.Use(authMiddleware.RequireRole("ADMIN"))
			{
				adminUsers.POST("", handlers.User.CreateUser)
				adminUsers.POST("/invite", handlers.User.CreateUserInvitation)
				adminUsers.PATCH("/:id", handlers.User.UpdateUser)
				adminUsers.DELETE("/:id", handlers.User.DeleteUser)
			}

			// Drivers routes (Admin and Company)
			drivers := protected.Group("/drivers")
			drivers.Use(authMiddleware.RequireRole("ADMIN", "COMPANY"))
			{
				drivers.GET("", handlers.Driver.ListDrivers)
				drivers.POST("", handlers.Driver.CreateDriver)
				drivers.GET("/:id", handlers.Driver.GetDriver)
				drivers.PATCH("/:id", handlers.Driver.UpdateDriver)
				drivers.DELETE("/:id", handlers.Driver.DeleteDriver)
				drivers.GET("/:id/kpis", handlers.Driver.GetDriverKPIs)
			}

			// Reservations routes (All authenticated users)
			reservations := protected.Group("/reservations")
			{
				reservations.GET("", handlers.Reservation.ListReservations)
				reservations.GET("/my", handlers.Reservation.GetMyReservations) // User's own reservations
				reservations.POST("", handlers.Reservation.CreateReservation)
				// Specific routes MUST come before generic /:id routes
				reservations.PATCH("/:id/status", handlers.Reservation.ChangeReservationStatus)
				reservations.PATCH("/:id/driver", handlers.Reservation.AssignDriver)
				reservations.GET("/:id/test", func(c *gin.Context) { c.JSON(200, gin.H{"test": "works", "id": c.Param("id")}) })
				reservations.GET("/:id/timeline", handlers.Reservation.GetReservationTimeline)
				reservations.POST("/:id/timeline", handlers.Reservation.AddTimelineEvent)
				// Generic routes come last
				reservations.GET("/:id", handlers.Reservation.GetReservation)
				reservations.PATCH("/:id", handlers.Reservation.UpdateReservation)
			}

			// Payments routes (All authenticated users)
			payments := protected.Group("/payments")
			{
				payments.POST("", handlers.Payment.CreatePayment)
				payments.GET("/:id", handlers.Payment.GetPayment)
				payments.POST("/:id/simulate", handlers.Payment.SimulatePayment)
			}

			// Companies routes (Admin can see all, User and Company can see their org)
			companies := protected.Group("/companies")
			companies.Use(authMiddleware.RequireRole("ADMIN", "USER", "COMPANY"))
			companies.Use(authMiddleware.RequireOrgScope())
			{
				companies.GET("", handlers.Company.ListCompanies)
				companies.GET("/:id", handlers.Company.GetCompany)
			}

			// Admin-only company operations (separate group to avoid middleware conflicts)
			adminCompanies := protected.Group("/companies")
			adminCompanies.Use(authMiddleware.RequireRole("ADMIN"))
			{
				adminCompanies.POST("", handlers.Company.CreateCompany)
				adminCompanies.PUT("/:id", handlers.Company.UpdateCompany)
				adminCompanies.DELETE("/:id", handlers.Company.DeleteCompany)
			}

			// Vehicles routes (Admin and Company)
			if handlers.Vehicle != nil {
				vehicles := protected.Group("/vehicles")
				vehicles.Use(authMiddleware.RequireRole("ADMIN", "COMPANY"))
				{
					vehicles.GET("", handlers.Vehicle.ListVehicles)
					vehicles.GET("/:id", handlers.Vehicle.GetVehicle)
					vehicles.POST("", handlers.Vehicle.CreateVehicle)
					vehicles.PUT("/:id", handlers.Vehicle.UpdateVehicle)
					vehicles.DELETE("/:id", handlers.Vehicle.DeleteVehicle)
					vehicles.POST("/:id/assign", handlers.Vehicle.AssignVehicleToDriver)
				}
			}

			// Support routes (All authenticated users)
			if handlers.Support != nil {
				support := protected.Group("/support")
				{
					support.POST("/contact", handlers.Support.ContactSupport)
				}
			}

			// Company user management routes (Company admins only)
			companyUsers := protected.Group("/company/users")
			companyUsers.Use(authMiddleware.RequireRole("COMPANY"))
			{
				companyUsers.POST("/invite", handlers.User.InviteCompanyUser)
			}
		}
	}
}
