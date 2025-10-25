// @title Turivo Backend API
// @version 1.0
// @description API backend for Turivo transportation management system
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.turivo.com/support
// @contact.email support@turivo.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	_ "turivo-backend/docs"
	"turivo-backend/internal/infrastructure/auth"
	"turivo-backend/internal/infrastructure/config"
	"turivo-backend/internal/infrastructure/email"
	"turivo-backend/internal/infrastructure/logging"
	"turivo-backend/internal/infrastructure/payment"
	"turivo-backend/internal/infrastructure/repository"
	"turivo-backend/internal/interface/http/handler"
	"turivo-backend/internal/interface/http/middleware"
	"turivo-backend/internal/interface/http/router"
	"turivo-backend/internal/usecase"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	// Initialize logger
	logger, err := logging.New(cfg.Log.Level)
	if err != nil {
		panic(fmt.Sprintf("Failed to initialize logger: %v", err))
	}
	defer logger.Sync()

	logger.Info("Starting Turivo Backend API",
		zap.String("port", cfg.HTTP.Port),
		zap.String("log_level", cfg.Log.Level),
	)

	// Connect to database
	dbPool, err := pgxpool.New(context.Background(), cfg.DB.DSN)
	if err != nil {
		logger.Fatal("Failed to connect to database", zap.Error(err))
	}
	defer dbPool.Close()

	// Test database connection
	if err := dbPool.Ping(context.Background()); err != nil {
		logger.Fatal("Failed to ping database", zap.Error(err))
	}
	logger.Info("Database connection established")

	// Create SQL DB for repositories that need it
	sqlDB, err := sql.Open("pgx", cfg.DB.DSN)
	if err != nil {
		logger.Fatal("Failed to create SQL DB connection", zap.Error(err))
	}
	defer sqlDB.Close()

	// Initialize validator
	validate := validator.New()

	// Initialize services
	passwordService := auth.NewPasswordService()
	paymentGateway := payment.NewWebpayMockGateway(logger)
	emailService := email.NewSMTPService(email.SMTPConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	}, logger)

	// Initialize repositories
	userRepo := repository.NewUserRepository(dbPool)
	refreshTokenRepo := repository.NewRefreshTokenRepository(dbPool)
	driverRepo := repository.NewDriverRepository(dbPool)
	reservationRepo := repository.NewReservationRepository(dbPool, logger)
	paymentRepo := repository.NewPaymentRepository(dbPool)
	registrationTokenRepo := repository.NewRegistrationTokenRepository(sqlDB, logger)
	companyRepo := repository.NewCompanyRepository(sqlDB, logger)
	vehicleRepo := repository.NewVehicleRepository(sqlDB, logger)

	// Initialize use cases
	authUseCase := usecase.NewAuthUseCase(userRepo, refreshTokenRepo, passwordService, cfg.JWT.Secret, cfg.JWT.AccessTTL, cfg.JWT.RefreshTTL, logger)
	userUseCase := usecase.NewUserUseCase(userRepo, registrationTokenRepo, passwordService, emailService, logger)
	driverUseCase := usecase.NewDriverUseCase(driverRepo, userRepo, emailService, registrationTokenRepo, passwordService, logger)
	reservationUseCase := usecase.NewReservationUseCase(reservationRepo, driverRepo, userRepo, emailService, logger)
	paymentUseCase := usecase.NewPaymentUseCase(paymentRepo, reservationRepo, paymentGateway, logger)
	companyUseCase := usecase.NewCompanyUseCase(companyRepo, logger)
	vehicleUseCase := usecase.NewVehicleUseCase(vehicleRepo, driverRepo, logger)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authUseCase, validate, logger)
	userHandler := handler.NewUserHandler(userUseCase, validate, logger)
	driverHandler := handler.NewDriverHandler(driverUseCase, validate, logger)
	driverDashboardHandler := handler.NewDriverDashboardHandler(driverUseCase, logger)
	reservationHandler := handler.NewReservationHandler(reservationUseCase, validate, logger)
	paymentHandler := handler.NewPaymentHandler(paymentUseCase, validate, logger)
	companyHandler := handler.NewCompanyHandler(companyUseCase, validate, logger)
	vehicleHandler := handler.NewVehicleHandler(vehicleUseCase, validate, logger)
	supportHandler := handler.NewSupportHandler(emailService, userRepo, validate, logger)
	billingHandler := handler.NewBillingHandler(paymentUseCase, logger)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authUseCase, logger)

	// Set Gin mode based on log level
	if cfg.Log.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize Gin router
	ginRouter := gin.New()

	// Add middlewares
	ginRouter.Use(middleware.RequestID())
	ginRouter.Use(middleware.Logger(logger))
	ginRouter.Use(middleware.Recovery(logger))
	ginRouter.Use(middleware.CORS(cfg.CORS.Origins)) // Use configured CORS origins

	// Health check endpoint
	ginRouter.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "Turivo Backend API is running",
			"version": "1.0.0",
		})
	})

	// Swagger documentation
	ginRouter.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Create admin handler
	adminHandler := handler.NewAdminHandler(
		*userRepo,
		*driverRepo,
		*vehicleRepo,
		*reservationRepo,
		*companyRepo,
		logger,
	)

	// Create company detail handler
	companyDetailHandler := handler.NewCompanyDetailHandler(
		*companyRepo,
		*userRepo,
		logger,
	)

	// Setup API routes
	router.SetupRoutes(ginRouter, router.RouteHandlers{
		Auth:            authHandler,
		User:            userHandler,
		Driver:          driverHandler,
		DriverDashboard: driverDashboardHandler,
		Reservation:     reservationHandler,
		Payment:         paymentHandler,
		Company:         companyHandler,
		CompanyDetail:   companyDetailHandler,
		Vehicle:         vehicleHandler,
		Support:         supportHandler,
		Admin:           adminHandler,
		Billing:         billingHandler,
	}, authMiddleware)

	// Start server
	port := ":" + cfg.HTTP.Port
	logger.Info("Server starting", zap.String("address", port))
	if err := ginRouter.Run(port); err != nil {
		logger.Fatal("Failed to start server", zap.Error(err))
	}
}
