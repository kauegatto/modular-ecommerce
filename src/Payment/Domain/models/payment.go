package models

import (
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type Money int64
type PaymentID = uuid.UUID
type PaymentStatus string
type PaymentKind string

const (
	PaymentStatusPlaced        PaymentStatus = "PAYMENT_CREATED"
	PaymentStatusPending       PaymentStatus = "PAYMENT_PENDING"
	PaymentStatusCompleted     PaymentStatus = "PAYMENT_COMPLETED"
	PaymentStatusPendingRefund PaymentStatus = "PAYMENT_PENDING_REFUND"
	PaymentStatusRefunded      PaymentStatus = "PAYMENT_REFUNDED"
	PaymentKindDebit           PaymentKind   = "DEBIT"
	PaymentKindCredit          PaymentKind   = "CREDIT"
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

func (p *Payment) AddExternalIntegratorID(id string) error {
	p.ExternalIntegratorID = id
	return nil
}

func (p *Payment) PendingPayment() error {
	if p.Status != PaymentStatusPlaced {
		return fmt.Errorf("cannot set Payment to pending with status: %s", p.Status)
	}
	p.Status = PaymentStatusPending
	return nil
}

func (p *Payment) RequestRefund() error {
	if p.Status == PaymentStatusPendingRefund {
		return fmt.Errorf("Payment is already pending refund")
	}
	if p.Status == PaymentStatusRefunded {
		return fmt.Errorf("Payment was alreday refunded")
	}

	p.Status = PaymentStatusPendingRefund
	return nil
}

func (p *Payment) ConfirmRefund() error {
	if p.Status == PaymentStatusRefunded {
		return fmt.Errorf("Payment was alreday refunded")
	}

	p.Status = PaymentStatusRefunded
	return nil
}

func NewPayment(orderId string, totalAmount Money) (*Payment, error) {
	if totalAmount < 100 {
		return nil, fmt.Errorf("can not create order less than R$1BRL")
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
	if kind == "" || status == "" {
		log.Fatalf("treat kind: %s status %s", kind, status)
	}
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
