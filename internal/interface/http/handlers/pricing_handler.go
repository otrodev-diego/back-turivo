package handlers

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
	"turivo-backend/internal/usecase"
)

type PricingHandler struct {
	pricingUseCase *usecase.PricingUseCase
	logger         *zap.Logger
}

func NewPricingHandler(pricingUseCase *usecase.PricingUseCase, logger *zap.Logger) *PricingHandler {
	return &PricingHandler{
		pricingUseCase: pricingUseCase,
		logger:         logger,
	}
}

// PricingQuoteRequest representa la entrada del endpoint de cotización
type PricingQuoteRequest struct {
	ServiceCode   string   `json:"serviceCode" binding:"required"`
	DistanceKm    *float64 `json:"distanceKm,omitempty"`
	VehicleTypeID string   `json:"vehicleTypeId,omitempty"`
	SegmentID     string   `json:"segmentId,omitempty"`
	ZoneID        string   `json:"zoneId" binding:"required"`
	ScheduleID    string   `json:"scheduleId" binding:"required"`
	CurrencyCode  string   `json:"currencyCode" binding:"required"`
	Paradas       *int     `json:"paradas,omitempty"`
	HorasEspera   *float64 `json:"horasEspera,omitempty"`
}

// PricingQuoteResponse representa la salida del endpoint de cotización
type PricingQuoteResponse struct {
	ServiceCode  string             `json:"serviceCode"`
	Mode         string             `json:"mode"`
	Currency     string             `json:"currency"`
	FinalFare    float64            `json:"finalFare"`
	Commission   float64            `json:"commission"`
	DriverPayout float64            `json:"driverPayout"`
	Inputs       map[string]any     `json:"inputs"`
	Breakdown    map[string]float64 `json:"breakdown"`
}

// Quote maneja la cotización de precios
func (h *PricingHandler) Quote(c *gin.Context) {
	var req PricingQuoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Error binding request", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
		return
	}

	h.logger.Info("Processing pricing quote", zap.String("serviceCode", req.ServiceCode))

	// Validar campos requeridos según el tipo de servicio
	if err := h.validateRequest(c.Request.Context(), &req); err != nil {
		h.logger.Error("Validation error", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Calcular precio usando el use case
	result, err := h.pricingUseCase.CalculatePrice(c.Request.Context(), &domain.PricingRequest{
		ServiceCode:   req.ServiceCode,
		DistanceKm:    req.DistanceKm,
		VehicleTypeID: req.VehicleTypeID,
		SegmentID:     req.SegmentID,
		ZoneID:        req.ZoneID,
		ScheduleID:    req.ScheduleID,
		CurrencyCode:  req.CurrencyCode,
		Paradas:       req.Paradas,
		HorasEspera:   req.HorasEspera,
	})

	if err != nil {
		h.logger.Error("Error calculating price", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error calculating price"})
		return
	}

	// Construir respuesta
	response := PricingQuoteResponse{
		ServiceCode:  result.ServiceCode,
		Mode:         result.Mode,
		Currency:     result.Currency,
		FinalFare:    result.FinalFare,
		Commission:   result.Commission,
		DriverPayout: result.DriverPayout,
		Inputs: map[string]any{
			"distanceKm":    req.DistanceKm,
			"vehicleTypeId": req.VehicleTypeID,
			"segmentId":     req.SegmentID,
			"zoneId":        req.ZoneID,
			"scheduleId":    req.ScheduleID,
		},
		Breakdown: result.Breakdown,
	}

	h.logger.Info("Pricing quote completed",
		zap.Float64("finalFare", result.FinalFare),
		zap.Float64("commission", result.Commission),
		zap.Float64("driverPayout", result.DriverPayout))

	c.JSON(http.StatusOK, response)
}

// validateRequest valida la entrada según el tipo de servicio
func (h *PricingHandler) validateRequest(ctx context.Context, req *PricingQuoteRequest) error {
	// Validar que el serviceCode existe y está activo
	service, err := h.pricingUseCase.GetServiceByCode(ctx, req.ServiceCode)
	if err != nil {
		return err
	}

	// Validar campos según el modo de servicio
	if service.Mode == "transfer" {
		if req.DistanceKm == nil || *req.DistanceKm <= 0 {
			return domain.ErrInvalidInput
		}
		if req.VehicleTypeID == "" {
			return domain.ErrInvalidInput
		}
		if req.SegmentID == "" {
			return domain.ErrInvalidInput
		}
	}

	// Validar límites razonables
	if req.DistanceKm != nil && *req.DistanceKm > 1000 {
		return domain.ErrInvalidInput
	}

	return nil
}
