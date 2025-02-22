package outgoing

import (
	"time"

	"github.com/google/uuid"
)

type PaymentCompleted struct {
	OrderID   string
	PaymentID uuid.UUID
	Amount    string
	Time      time.Time
}

func (p PaymentCompleted) Name() string {
	return "PaymentCompleted"
}

type PaymentCreated struct {
	OrderID   string
	PaymentID uuid.UUID
	Amount    string
	Time      time.Time
}

func (p PaymentCreated) Name() string {
	return "PaymentCreated"
}

type PaymentRefundRequested struct {
	OrderID   string
	PaymentID uuid.UUID
	Amount    string
	Time      time.Time
}

func (p PaymentRefundRequested) Name() string {
	return "PaymentCancelRequested"
}

type PaymentCancelConfirmed struct {
	OrderID   string
	PaymentID uuid.UUID
	Amount    string
	Time      time.Time
}

func (p PaymentCancelConfirmed) Name() string {
	return "PaymentCancelConfirmed"
}

type PaymentCancelDenied struct {
	OrderID   string
	PaymentID uuid.UUID
	Amount    string
	Time      time.Time
}

func (p PaymentCancelDenied) Name() string {
	return "PaymentCancelDenied"
}
