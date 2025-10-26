package usecase

import (
	"context"
	"fmt"
	"math"

	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type PricingUseCase struct {
	pricingRepo domain.PricingRepository
	logger      *zap.Logger
}

func NewPricingUseCase(pricingRepo domain.PricingRepository, logger *zap.Logger) *PricingUseCase {
	return &PricingUseCase{
		pricingRepo: pricingRepo,
		logger:      logger,
	}
}

// CalculatePrice calcula el precio según el algoritmo de Turivo
func (uc *PricingUseCase) CalculatePrice(ctx context.Context, req *domain.PricingRequest) (*domain.PricingResult, error) {
	uc.logger.Info("Calculating price", zap.String("serviceCode", req.ServiceCode))

	// 1. Obtener configuración global
	settings, err := uc.pricingRepo.GetSettings(ctx)
	if err != nil {
		return nil, err
	}

	// 2. Obtener servicio
	service, err := uc.pricingRepo.GetServiceByCode(ctx, req.ServiceCode)
	if err != nil {
		return nil, err
	}

	// 3. Obtener factores
	factors, err := uc.getFactors(ctx, req)
	if err != nil {
		return nil, err
	}

	// 4. Calcular precio según el modo
	var finalFare float64
	var breakdown map[string]float64

	if service.Mode == "transfer" {
		finalFare, breakdown = uc.calculateTransfer(req, service, factors, settings)
	} else if service.Mode == "tour" {
		finalFare, breakdown = uc.calculateTour(req, service, factors, settings)
	} else {
		return nil, fmt.Errorf("unsupported service mode: %s", service.Mode)
	}

	// 5. Aplicar redondeo
	finalFare = uc.roundPrice(finalFare, settings.RoundingDecimals)
	commission := uc.roundPrice(finalFare*settings.CommissionRate, settings.RoundingDecimals)
	driverPayout := uc.roundPrice(finalFare-commission, settings.RoundingDecimals)

	// 6. Obtener tasa de cambio si es necesario
	exchangeRate := 1.0
	if req.CurrencyCode != "CLP" {
		rate, err := uc.pricingRepo.GetCurrencyRate(ctx, req.CurrencyCode)
		if err != nil {
			return nil, err
		}
		exchangeRate = rate
	}

	// 7. Convertir a moneda solicitada
	finalFare = finalFare / exchangeRate
	commission = commission / exchangeRate
	driverPayout = driverPayout / exchangeRate

	result := &domain.PricingResult{
		ServiceCode:  req.ServiceCode,
		Mode:         service.Mode,
		Currency:     req.CurrencyCode,
		FinalFare:    finalFare,
		Commission:   commission,
		DriverPayout: driverPayout,
		Breakdown:    breakdown,
	}

	uc.logger.Info("Price calculated successfully",
		zap.Float64("finalFare", finalFare),
		zap.Float64("commission", commission),
		zap.Float64("driverPayout", driverPayout))

	return result, nil
}

// calculateTransfer calcula precio para transfers
func (uc *PricingUseCase) calculateTransfer(req *domain.PricingRequest, service *domain.PricingService, factors *domain.PricingFactors, settings *domain.PricingSettings) (float64, map[string]float64) {
	// Fórmula: product = (base_per_km * distancia) * Fv * Fs * Fz * Fh
	base := settings.BasePerKmCLP * *req.DistanceKm
	product := base * factors.Vehicle * factors.Segment * factors.Zone * factors.Schedule

	// Log detallado para debug
	uc.logger.Info("Transfer calculation details",
		zap.Float64("base", base),
		zap.Float64("product", product),
		zap.Float64("minFare", service.MinFareCLP),
		zap.Float64("vehicleFactor", factors.Vehicle),
		zap.Float64("segmentFactor", factors.Segment),
		zap.Float64("zoneFactor", factors.Zone),
		zap.Float64("scheduleFactor", factors.Schedule))

	// Aplicar tarifa mínima
	finalFare := math.Max(product, service.MinFareCLP)

	// Agregar extras para rutas integradas
	if req.Paradas != nil && *req.Paradas > 0 {
		finalFare += float64(*req.Paradas) * 3000
	}
	if req.HorasEspera != nil && *req.HorasEspera > 0 {
		finalFare += *req.HorasEspera * 16000
	}

	breakdown := map[string]float64{
		"basePerKmCLP":   settings.BasePerKmCLP,
		"factorVehicle":  factors.Vehicle,
		"factorSegment":  factors.Segment,
		"factorZone":     factors.Zone,
		"factorSchedule": factors.Schedule,
		"minFareCLP":     service.MinFareCLP,
	}

	return finalFare, breakdown
}

// calculateTour calcula precio para tours
func (uc *PricingUseCase) calculateTour(req *domain.PricingRequest, service *domain.PricingService, factors *domain.PricingFactors, settings *domain.PricingSettings) (float64, map[string]float64) {
	// Fórmula: base = max(base_flat_clp, min_fare_clp) * Fz * Fh
	base := math.Max(service.BaseFlatCLP, service.MinFareCLP)
	finalFare := base * factors.Zone * factors.Schedule

	breakdown := map[string]float64{
		"baseFlatCLP":    service.BaseFlatCLP,
		"minFareCLP":     service.MinFareCLP,
		"factorZone":     factors.Zone,
		"factorSchedule": factors.Schedule,
	}

	return finalFare, breakdown
}

// getFactors obtiene todos los factores necesarios
func (uc *PricingUseCase) getFactors(ctx context.Context, req *domain.PricingRequest) (*domain.PricingFactors, error) {
	factors := &domain.PricingFactors{
		Vehicle:  1.0,
		Segment:  1.0,
		Zone:     1.0,
		Schedule: 1.0,
	}

	// Obtener factor de zona
	zoneFactor, err := uc.pricingRepo.GetZoneFactor(ctx, req.ZoneID)
	if err != nil {
		return nil, err
	}
	factors.Zone = zoneFactor

	// Obtener factor de horario (específico por servicio)
	scheduleFactor, err := uc.pricingRepo.GetScheduleFactorByService(ctx, req.ScheduleID, req.ServiceCode)
	if err != nil {
		return nil, err
	}
	factors.Schedule = scheduleFactor

	// Para transfers, obtener factores de vehículo y segmento
	if req.VehicleTypeID != "" {
		vehicleFactor, err := uc.pricingRepo.GetVehicleFactor(ctx, req.VehicleTypeID)
		if err != nil {
			return nil, err
		}
		factors.Vehicle = vehicleFactor
	}

	if req.SegmentID != "" {
		segmentFactor, err := uc.pricingRepo.GetSegmentFactor(ctx, req.SegmentID)
		if err != nil {
			return nil, err
		}
		factors.Segment = segmentFactor
	}

	return factors, nil
}

// roundPrice redondea el precio según la configuración
func (uc *PricingUseCase) roundPrice(price float64, decimals int) float64 {
	multiplier := math.Pow(10, float64(decimals))
	return math.Round(price*multiplier) / multiplier
}

// GetServiceByCode obtiene un servicio por código
func (uc *PricingUseCase) GetServiceByCode(ctx context.Context, code string) (*domain.PricingService, error) {
	return uc.pricingRepo.GetServiceByCode(ctx, code)
}
