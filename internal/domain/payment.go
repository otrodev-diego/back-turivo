package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type PaymentGateway string

const (
	PaymentGatewayWebpayPlus PaymentGateway = "WEBPAY_PLUS"
)

type PaymentStatus string

const (
	PaymentStatusApproved PaymentStatus = "APPROVED"
	PaymentStatusRejected PaymentStatus = "REJECTED"
	PaymentStatusPending  PaymentStatus = "PENDING"
)

type Payment struct {
	ID             uuid.UUID              `json:"id"`
	ReservationID  string                 `json:"reservation_id"`
	Gateway        PaymentGateway         `json:"gateway"`
	Amount         float64                `json:"amount"`
	Currency       string                 `json:"currency"`
	Status         PaymentStatus          `json:"status"`
	TransactionRef *string                `json:"transaction_ref,omitempty"`
	Payload        map[string]interface{} `json:"payload,omitempty"`
	CreatedAt      time.Time              `json:"created_at"`

	// Related data
	Reservation *Reservation `json:"reservation,omitempty"`
}

type CreatePaymentRequest struct {
	ReservationID string         `json:"reservation_id" validate:"required"`
	Method        PaymentGateway `json:"method" validate:"required"`
}

type SimulatePaymentRequest struct {
	Result PaymentStatus `json:"result" validate:"required,oneof=APPROVED REJECTED"`
}

// Helper methods for JSON marshaling/unmarshaling
func (p *Payment) MarshalPayload() ([]byte, error) {
	if p.Payload == nil {
		return []byte("{}"), nil
	}
	return json.Marshal(p.Payload)
}

func (p *Payment) UnmarshalPayload(data []byte) error {
	if len(data) == 0 {
		p.Payload = make(map[string]interface{})
		return nil
	}
	return json.Unmarshal(data, &p.Payload)
}

type PaymentRepository interface {
	Create(payment *Payment) error
	GetByID(id uuid.UUID) (*Payment, error)
	GetByReservationID(reservationID string) ([]*Payment, error)
	Update(id uuid.UUID, status PaymentStatus, transactionRef *string, payload map[string]interface{}) (*Payment, error)
}

// PaymentGatewayService interface for payment processing
type PaymentGatewayService interface {
	ProcessPayment(payment *Payment) (*PaymentResult, error)
	SimulatePayment(paymentID uuid.UUID, result PaymentStatus) (*PaymentResult, error)
}

type PaymentResult struct {
	PaymentID      uuid.UUID              `json:"payment_id"`
	Status         PaymentStatus          `json:"status"`
	TransactionRef *string                `json:"transaction_ref,omitempty"`
	Payload        map[string]interface{} `json:"payload,omitempty"`
	Message        string                 `json:"message"`
}
