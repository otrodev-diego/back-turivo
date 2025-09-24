package payment

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type WebpayMockGateway struct {
	logger *zap.Logger
}

func NewWebpayMockGateway(logger *zap.Logger) domain.PaymentGatewayService {
	return &WebpayMockGateway{
		logger: logger,
	}
}

func (w *WebpayMockGateway) ProcessPayment(payment *domain.Payment) (*domain.PaymentResult, error) {
	w.logger.Info("Processing payment with Webpay Mock",
		zap.String("payment_id", payment.ID.String()),
		zap.Float64("amount", payment.Amount))

	// Simulate processing time
	time.Sleep(100 * time.Millisecond)

	// Mock transaction reference
	transactionRef := fmt.Sprintf("WP_%d_%s", time.Now().Unix(), payment.ID.String()[:8])

	// Mock payload with typical Webpay response structure
	payload := map[string]interface{}{
		"vci":        "TSY",
		"amount":     payment.Amount,
		"status":     "AUTHORIZED",
		"buy_order":  payment.ReservationID,
		"session_id": payment.ID.String(),
		"card_detail": map[string]interface{}{
			"card_number": "XXXX-XXXX-XXXX-1234",
			"card_type":   "Visa",
		},
		"accounting_date":     time.Now().Format("0102"),
		"transaction_date":    time.Now().Format(time.RFC3339),
		"authorization_code":  fmt.Sprintf("AUTH%d", time.Now().Unix()%100000),
		"payment_type_code":   "VN",
		"response_code":       0,
		"installments_amount": payment.Amount,
		"installments_number": 1,
		"balance":             0.0,
	}

	return &domain.PaymentResult{
		PaymentID:      payment.ID,
		Status:         domain.PaymentStatusApproved,
		TransactionRef: &transactionRef,
		Payload:        payload,
		Message:        "Payment processed successfully",
	}, nil
}

func (w *WebpayMockGateway) SimulatePayment(paymentID uuid.UUID, result domain.PaymentStatus) (*domain.PaymentResult, error) {
	w.logger.Info("Simulating payment result",
		zap.String("payment_id", paymentID.String()),
		zap.String("result", string(result)))

	var transactionRef *string
	var payload map[string]interface{}
	var message string

	switch result {
	case domain.PaymentStatusApproved:
		ref := fmt.Sprintf("WP_SIM_%d_%s", time.Now().Unix(), paymentID.String()[:8])
		transactionRef = &ref
		payload = map[string]interface{}{
			"vci":        "TSY",
			"status":     "AUTHORIZED",
			"session_id": paymentID.String(),
			"card_detail": map[string]interface{}{
				"card_number": "XXXX-XXXX-XXXX-1234",
				"card_type":   "Visa",
			},
			"accounting_date":     time.Now().Format("0102"),
			"transaction_date":    time.Now().Format(time.RFC3339),
			"authorization_code":  fmt.Sprintf("AUTH%d", time.Now().Unix()%100000),
			"payment_type_code":   "VN",
			"response_code":       0,
			"installments_number": 1,
			"balance":             0.0,
		}
		message = "Payment approved successfully"

	case domain.PaymentStatusRejected:
		payload = map[string]interface{}{
			"vci":              "TSN",
			"status":           "REJECTED",
			"session_id":       paymentID.String(),
			"response_code":    -1,
			"rejection_reason": "Insufficient funds",
			"transaction_date": time.Now().Format(time.RFC3339),
		}
		message = "Payment was rejected"

	default:
		payload = map[string]interface{}{
			"status":           "PENDING",
			"session_id":       paymentID.String(),
			"transaction_date": time.Now().Format(time.RFC3339),
		}
		message = "Payment is pending"
	}

	return &domain.PaymentResult{
		PaymentID:      paymentID,
		Status:         result,
		TransactionRef: transactionRef,
		Payload:        payload,
		Message:        message,
	}, nil
}
