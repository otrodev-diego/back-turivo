package usecase

import (
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type PaymentUseCase struct {
	paymentRepo     domain.PaymentRepository
	reservationRepo domain.ReservationRepository
	paymentGateway  domain.PaymentGatewayService
	logger          *zap.Logger
}

func NewPaymentUseCase(
	paymentRepo domain.PaymentRepository,
	reservationRepo domain.ReservationRepository,
	paymentGateway domain.PaymentGatewayService,
	logger *zap.Logger,
) *PaymentUseCase {
	return &PaymentUseCase{
		paymentRepo:     paymentRepo,
		reservationRepo: reservationRepo,
		paymentGateway:  paymentGateway,
		logger:          logger,
	}
}

func (uc *PaymentUseCase) CreatePayment(req domain.CreatePaymentRequest) (*domain.Payment, error) {
	uc.logger.Info("Creating payment", zap.String("reservation_id", req.ReservationID))

	// Check if reservation exists
	reservation, err := uc.reservationRepo.GetByID(req.ReservationID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrReservationNotFound
		}
		uc.logger.Error("Failed to get reservation for payment", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Check if reservation already has approved payment
	existingPayments, err := uc.paymentRepo.GetByReservationID(req.ReservationID)
	if err != nil {
		uc.logger.Error("Failed to check existing payments", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	for _, payment := range existingPayments {
		if payment.Status == domain.PaymentStatusApproved {
			return nil, domain.ErrPaymentAlreadyPaid
		}
	}

	// Get amount from reservation
	if reservation.Amount == nil {
		uc.logger.Error("Reservation has no amount set", zap.String("reservation_id", req.ReservationID))
		return nil, domain.ErrInvalidInput
	}

	// Create payment
	payment := &domain.Payment{
		ID:            uuid.New(),
		ReservationID: req.ReservationID,
		Gateway:       req.Method,
		Amount:        *reservation.Amount,
		Currency:      "CLP",
		Status:        domain.PaymentStatusPending,
		CreatedAt:     time.Now(),
	}

	if err := uc.paymentRepo.Create(payment); err != nil {
		uc.logger.Error("Failed to create payment", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Process payment through gateway
	result, err := uc.paymentGateway.ProcessPayment(payment)
	if err != nil {
		uc.logger.Error("Failed to process payment", zap.Error(err))
		return nil, domain.ErrPaymentFailed
	}

	// Update payment with result
	updatedPayment, err := uc.paymentRepo.Update(
		payment.ID,
		result.Status,
		result.TransactionRef,
		result.Payload,
	)
	if err != nil {
		uc.logger.Error("Failed to update payment with result", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	uc.logger.Info("Payment created and processed",
		zap.String("payment_id", payment.ID.String()),
		zap.String("status", string(result.Status)))

	return updatedPayment, nil
}

func (uc *PaymentUseCase) GetPaymentByID(id uuid.UUID) (*domain.Payment, error) {
	payment, err := uc.paymentRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrPaymentNotFound
		}
		uc.logger.Error("Failed to get payment by ID", zap.Error(err), zap.String("payment_id", id.String()))
		return nil, domain.ErrInternalError
	}

	return payment, nil
}

func (uc *PaymentUseCase) SimulatePayment(paymentID uuid.UUID, result domain.PaymentStatus) (*domain.Payment, error) {
	uc.logger.Info("Simulating payment",
		zap.String("payment_id", paymentID.String()),
		zap.String("result", string(result)))

	// Check if payment exists
	payment, err := uc.paymentRepo.GetByID(paymentID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrPaymentNotFound
		}
		uc.logger.Error("Failed to get payment for simulation", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Check if payment is already processed
	if payment.Status != domain.PaymentStatusPending {
		return nil, domain.ErrPaymentAlreadyPaid
	}

	// Simulate payment through gateway
	simulationResult, err := uc.paymentGateway.SimulatePayment(paymentID, result)
	if err != nil {
		uc.logger.Error("Failed to simulate payment", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Update payment with simulation result
	updatedPayment, err := uc.paymentRepo.Update(
		paymentID,
		simulationResult.Status,
		simulationResult.TransactionRef,
		simulationResult.Payload,
	)
	if err != nil {
		uc.logger.Error("Failed to update payment with simulation result", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	uc.logger.Info("Payment simulation completed",
		zap.String("payment_id", paymentID.String()),
		zap.String("status", string(simulationResult.Status)))

	return updatedPayment, nil
}

func (uc *PaymentUseCase) GetPaymentsByReservationID(reservationID string) ([]*domain.Payment, error) {
	payments, err := uc.paymentRepo.GetByReservationID(reservationID)
	if err != nil {
		uc.logger.Error("Failed to get payments by reservation ID", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	return payments, nil
}
