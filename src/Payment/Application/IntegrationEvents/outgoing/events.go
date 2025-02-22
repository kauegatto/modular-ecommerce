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

type PaymentRefundConfirmed struct {
	OrderID   string
	PaymentID uuid.UUID
	Amount    string
	Time      time.Time
}

func (p PaymentRefundConfirmed) Name() string {
	return "PaymentCancelConfirmed"
}

type PaymentRefundDenied struct {
	OrderID   string
	PaymentID uuid.UUID
	Amount    string
	Time      time.Time
}

func (p PaymentRefundDenied) Name() string {
	return "PaymentCancelDenied"
}
