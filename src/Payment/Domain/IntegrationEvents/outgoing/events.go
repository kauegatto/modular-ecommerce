package outgoing

import "time"

type PaymentCompletedV1 struct {
	OrderID   string    `json:"order_id"`
	Amount    float64   `json:"amount"`
	TimeStamp time.Time `json:"timestamp"`
}
