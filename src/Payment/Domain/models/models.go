package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Money int64
type PaymentID = uuid.UUID
type PaymentStatus string

const (
	PaymentStatusPlaced    PaymentStatus = "PAYMENT_CREATED"
	PaymentStatusPending   PaymentStatus = "PAYMENT_PENDING_PAYMENT"
	PaymentStatusCompleted PaymentStatus = "PAYMENT_COMPLETED"
	PaymentStatusCancelled PaymentStatus = "PAYMENT_CANCELLED"
)

type Payment struct {
	ID                   PaymentID
	OrderId              string
	ExternalIntegratorID string
	status               PaymentStatus `mapstructure:"status"`
	createdAt            time.Time     `mapstructure:"createdAt"`
	totalPrice           Money         `mapstructure:"totalPrice"`
}

func (p *Payment) CompletePayment() error {
	if p.status != PaymentStatusPlaced {
		return fmt.Errorf("cannot complete Payment with status: %s", p.status)
	}
	p.status = PaymentStatusCompleted
	return nil
}

func (p *Payment) PendingPayment() error {
	if p.status != PaymentStatusPlaced {
		return fmt.Errorf("cannot set Payment to pending with status: %s", p.status)
	}
	p.status = PaymentStatusPending
	return nil
}

func (p *Payment) CancelPayment() error {
	if p.status == PaymentStatusCancelled {
		return fmt.Errorf("Payment is already cancelled")
	}
	p.status = PaymentStatusCancelled
	return nil
}

func (p *Payment) TotalPrice() Money {
	return p.totalPrice
}

func (p *Payment) Status() PaymentStatus {
	return p.status
}

func (p *Payment) CreatedAt() time.Time {
	return p.createdAt
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
		status:               PaymentStatusPlaced,
		createdAt:            time.Now().UTC(),
		totalPrice:           totalAmount,
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
) *Payment {
	return &Payment{
		ID:                   ID,
		OrderId:              OrderId,
		ExternalIntegratorID: ExternalIntegratorID,
		status:               PaymentStatusPlaced,
		createdAt:            time.Now().UTC(),
		totalPrice:           TotalAmount,
	}
}
