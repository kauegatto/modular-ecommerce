package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Money int64
type PaymentID = uuid.UUID
type PaymentStatus string
type PaymentKind string

const (
	PaymentStatusPlaced    PaymentStatus = "PAYMENT_CREATED"
	PaymentStatusPending   PaymentStatus = "PAYMENT_PENDING_PAYMENT"
	PaymentStatusCompleted PaymentStatus = "PAYMENT_COMPLETED"
	PaymentStatusCancelled PaymentStatus = "PAYMENT_CANCELLED"
	PaymentKindDebit       PaymentKind   = "DEBIT"
	PaymentKindCredit      PaymentKind   = "CREDIT"
)

type Payment struct {
	ID                   PaymentID
	OrderId              string
	ExternalIntegratorID string
	Kind                 PaymentKind
	Status               PaymentStatus
	CreatedAt            time.Time
	TotalPrice           Money
}

func (p *Payment) CompletePayment() error {
	if p.Status != PaymentStatusPlaced {
		return fmt.Errorf("cannot complete Payment with status: %s", p.Status)
	}
	p.Status = PaymentStatusCompleted
	return nil
}

func (p *Payment) PendingPayment() error {
	if p.Status != PaymentStatusPlaced {
		return fmt.Errorf("cannot set Payment to pending with status: %s", p.Status)
	}
	p.Status = PaymentStatusPending
	return nil
}

func (p *Payment) CancelPayment() error {
	if p.Status == PaymentStatusCancelled {
		return fmt.Errorf("Payment is already cancelled")
	}
	p.Status = PaymentStatusCancelled
	return nil
}

func NewPayment(orderId string, totalAmount Money) (*Payment, error) {
	if totalAmount < 100 {
		return nil, fmt.Errorf("Can not create order less than R$1BRL")
	}

	id, _ := uuid.NewV7()
	Payment := &Payment{
		ID:                   id,
		OrderId:              orderId,
		ExternalIntegratorID: "",
		Status:               PaymentStatusPlaced,
		CreatedAt:            time.Now().UTC(),
		TotalPrice:           totalAmount,
	}
	return Payment, nil
}

func NewPaymentFromRehidration(
	ID PaymentID,
	OrderId string,
	ExternalIntegratorID string,
	status PaymentStatus,
	createdAt time.Time,
	TotalAmount Money,
	kind PaymentKind,
) *Payment {
	return &Payment{
		ID:                   ID,
		OrderId:              OrderId,
		ExternalIntegratorID: ExternalIntegratorID,
		Status:               status,
		Kind:                 kind,
		CreatedAt:            time.Now().UTC(),
		TotalPrice:           TotalAmount,
	}
}
