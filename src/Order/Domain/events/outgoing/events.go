package outgoing

import "time"

type OrderPlaced struct {
	OrderID    string    `json:"order_id"`
	CustomerID string    `json:"customer_id"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

type OrderCancelled struct {
	OrderID     string
	CancelledAt time.Time
	Reason      string
}
