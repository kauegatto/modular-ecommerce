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

type PaymentCancelled struct {
	OrderID   string
	PaymentID uuid.UUID
	Amount    string
	Time      time.Time
}

func (p PaymentCancelled) Name() string {
	return "PaymentCancelled"
}
