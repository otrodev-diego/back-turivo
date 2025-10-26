package repository

import (
	"context"
	"database/sql"

	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type PricingRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewPricingRepository(db *sql.DB, logger *zap.Logger) *PricingRepository {
	return &PricingRepository{
		db:     db,
		logger: logger,
	}
}

// GetSettings obtiene la configuración global de pricing
func (r *PricingRepository) GetSettings(ctx context.Context) (*domain.PricingSettings, error) {
	// Por ahora, retornar configuración hardcodeada según el documento
	// En producción, esto vendría de la base de datos
	settings := &domain.PricingSettings{
		BasePerKmCLP:     1200,
		CommissionRate:   0.20,
		DefaultCurrency:  "CLP",
		RoundingDecimals: 2, // Cambiar a 2 decimales para evitar redondeo prematuro
	}

	r.logger.Info("Retrieved pricing settings", zap.Float64("basePerKmCLP", settings.BasePerKmCLP))
	return settings, nil
}

// GetServiceByCode obtiene un servicio por código
func (r *PricingRepository) GetServiceByCode(ctx context.Context, code string) (*domain.PricingService, error) {
	// Mapeo de servicios según el documento de validación
	services := map[string]*domain.PricingService{
		"T004": {
			Code:       "T004",
			Name:       "Traslado Aeropuerto",
			Mode:       "transfer",
			MinFareCLP: 42000,
			Status:     "active",
		},
		"T003": {
			Code:       "T003",
			Name:       "Traslado Urbano",
			Mode:       "transfer",
			MinFareCLP: 12000,
			Status:     "active",
		},
		"T014": {
			Code:        "T014",
			Name:        "Tour Viña del Mar",
			Mode:        "tour",
			MinFareCLP:  350000,
			BaseFlatCLP: 250000,
			Status:      "active",
		},
		"T015": {
			Code:        "T015",
			Name:        "Tour Cajón del Maipo",
			Mode:        "tour",
			MinFareCLP:  250000,
			BaseFlatCLP: 250000,
			Status:      "active",
		},
		"T009": {
			Code:       "T009",
			Name:       "Ruta Integrada",
			Mode:       "transfer",
			MinFareCLP: 36000,
			Status:     "active",
		},
	}

	service, exists := services[code]
	if !exists {
		r.logger.Error("Service not found", zap.String("code", code))
		return nil, domain.ErrServiceNotFound
	}

	r.logger.Info("Retrieved service", zap.String("code", code), zap.String("mode", service.Mode))
	return service, nil
}

// GetVehicleFactor obtiene el factor de vehículo
func (r *PricingRepository) GetVehicleFactor(ctx context.Context, vehicleTypeID string) (float64, error) {
	factors := map[string]float64{
		"van_estandar":     1.0,
		"van_premium":      1.4,
		"minibus_estandar": 1.4,
		"minibus_premium":  2.0,
		"bus_estandar":     2.0,
		"bus_premium":      2.5,
		"sedan_ejecutivo":  1.2,
		"suv_premium":      2.0,
	}

	factor, exists := factors[vehicleTypeID]
	if !exists {
		r.logger.Error("Vehicle factor not found", zap.String("vehicleTypeID", vehicleTypeID))
		return 0, domain.ErrInvalidFactors
	}

	r.logger.Info("Retrieved vehicle factor", zap.String("vehicleTypeID", vehicleTypeID), zap.Float64("factor", factor))
	return factor, nil
}

// GetSegmentFactor obtiene el factor de segmento
func (r *PricingRepository) GetSegmentFactor(ctx context.Context, segmentID string) (float64, error) {
	factors := map[string]float64{
		"B2C": 1.0,
		"B2B": 0.9,
	}

	factor, exists := factors[segmentID]
	if !exists {
		r.logger.Error("Segment factor not found", zap.String("segmentID", segmentID))
		return 0, domain.ErrInvalidFactors
	}

	r.logger.Info("Retrieved segment factor", zap.String("segmentID", segmentID), zap.Float64("factor", factor))
	return factor, nil
}

// GetZoneFactor obtiene el factor de zona
func (r *PricingRepository) GetZoneFactor(ctx context.Context, zoneID string) (float64, error) {
	factors := map[string]float64{
		"urbana":        1.0,
		"mixta":         1.1,
		"rural":         1.2,
		"interregional": 1.3,
	}

	factor, exists := factors[zoneID]
	if !exists {
		r.logger.Error("Zone factor not found", zap.String("zoneID", zoneID))
		return 0, domain.ErrInvalidFactors
	}

	r.logger.Info("Retrieved zone factor", zap.String("zoneID", zoneID), zap.Float64("factor", factor))
	return factor, nil
}

// GetScheduleFactor obtiene el factor de horario
func (r *PricingRepository) GetScheduleFactor(ctx context.Context, scheduleID string) (float64, error) {
	factors := map[string]float64{
		"normal":   1.0,
		"punta":    1.3, // Factor base para transfers
		"nocturno": 1.2,
	}

	factor, exists := factors[scheduleID]
	if !exists {
		r.logger.Error("Schedule factor not found", zap.String("scheduleID", scheduleID))
		return 0, domain.ErrInvalidFactors
	}

	r.logger.Info("Retrieved schedule factor", zap.String("scheduleID", scheduleID), zap.Float64("factor", factor))
	return factor, nil
}

// GetScheduleFactorByService obtiene el factor de horario específico por servicio
func (r *PricingRepository) GetScheduleFactorByService(ctx context.Context, scheduleID string, serviceCode string) (float64, error) {
	// Factores específicos por servicio según el documento de validación
	if serviceCode == "T015" && scheduleID == "punta" {
		// Para T015 (Tour Cajón del Maipo), "punta" debe ser 1.2, no 1.3
		factor := 1.2
		r.logger.Info("Retrieved schedule factor by service", zap.String("scheduleID", scheduleID), zap.String("serviceCode", serviceCode), zap.Float64("factor", factor))
		return factor, nil
	}

	// Para todos los demás casos, usar factor estándar
	return r.GetScheduleFactor(ctx, scheduleID)
}

// GetCurrencyRate obtiene la tasa de cambio de moneda
func (r *PricingRepository) GetCurrencyRate(ctx context.Context, currencyCode string) (float64, error) {
	rates := map[string]float64{
		"CLP": 1.0,
		"PEN": 235.0,
		"USD": 0.0011,
	}

	rate, exists := rates[currencyCode]
	if !exists {
		r.logger.Error("Currency rate not found", zap.String("currencyCode", currencyCode))
		return 0, domain.ErrInvalidCurrency
	}

	r.logger.Info("Retrieved currency rate", zap.String("currencyCode", currencyCode), zap.Float64("rate", rate))
	return rate, nil
}

// AuditQuote registra la cotización para auditoría
func (r *PricingRepository) AuditQuote(ctx context.Context, request *domain.PricingRequest, result *domain.PricingResult) error {
	// Por ahora, solo loggear. En producción, guardar en BD
	r.logger.Info("Auditing pricing quote",
		zap.String("serviceCode", request.ServiceCode),
		zap.Float64("finalFare", result.FinalFare),
		zap.Float64("commission", result.Commission),
		zap.Float64("driverPayout", result.DriverPayout))

	return nil
}
