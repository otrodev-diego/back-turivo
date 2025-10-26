package domain

import (
	"context"
	"errors"
)

// Errores específicos de pricing
var (
	ErrServiceNotFound    = errors.New("service not found")
	ErrInvalidServiceCode = errors.New("invalid service code")
	ErrInvalidFactors     = errors.New("invalid factors")
	ErrInvalidCurrency    = errors.New("invalid currency")
)

// PricingRequest representa la entrada para cálculo de precios
type PricingRequest struct {
	ServiceCode   string   `json:"serviceCode"`
	DistanceKm    *float64 `json:"distanceKm,omitempty"`
	VehicleTypeID string   `json:"vehicleTypeId,omitempty"`
	SegmentID     string   `json:"segmentId,omitempty"`
	ZoneID        string   `json:"zoneId"`
	ScheduleID    string   `json:"scheduleId"`
	CurrencyCode  string   `json:"currencyCode"`
	Paradas       *int     `json:"paradas,omitempty"`
	HorasEspera   *float64 `json:"horasEspera,omitempty"`
}

// PricingResult representa el resultado del cálculo de precios
type PricingResult struct {
	ServiceCode  string             `json:"serviceCode"`
	Mode         string             `json:"mode"`
	Currency     string             `json:"currency"`
	FinalFare    float64            `json:"finalFare"`
	Commission   float64            `json:"commission"`
	DriverPayout float64            `json:"driverPayout"`
	Breakdown    map[string]float64 `json:"breakdown"`
}

// PricingService representa un servicio de pricing
type PricingService struct {
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Mode        string  `json:"mode"`
	MinFareCLP  float64 `json:"minFareCLP"`
	BaseFlatCLP float64 `json:"baseFlatCLP,omitempty"`
	Status      string  `json:"status"`
}

// PricingSettings representa la configuración global de pricing
type PricingSettings struct {
	BasePerKmCLP     float64 `json:"basePerKmCLP"`
	CommissionRate   float64 `json:"commissionRate"`
	DefaultCurrency  string  `json:"defaultCurrency"`
	RoundingDecimals int     `json:"roundingDecimals"`
}

// PricingFactors representa los factores de precio
type PricingFactors struct {
	Vehicle  float64 `json:"vehicle"`
	Segment  float64 `json:"segment"`
	Zone     float64 `json:"zone"`
	Schedule float64 `json:"schedule"`
}

// PricingRepository define la interfaz para el repositorio de pricing
type PricingRepository interface {
	GetSettings(ctx context.Context) (*PricingSettings, error)
	GetServiceByCode(ctx context.Context, code string) (*PricingService, error)
	GetVehicleFactor(ctx context.Context, vehicleTypeID string) (float64, error)
	GetSegmentFactor(ctx context.Context, segmentID string) (float64, error)
	GetZoneFactor(ctx context.Context, zoneID string) (float64, error)
	GetScheduleFactor(ctx context.Context, scheduleID string) (float64, error)
	GetScheduleFactorByService(ctx context.Context, scheduleID string, serviceCode string) (float64, error)
	GetCurrencyRate(ctx context.Context, currencyCode string) (float64, error)
	AuditQuote(ctx context.Context, request *PricingRequest, result *PricingResult) error
}
