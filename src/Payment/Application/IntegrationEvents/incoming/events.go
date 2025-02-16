package incoming

import (
	"time"
)

type OrderPlaced struct {
	OrderID    string    `json:"order_id"`
	CustomerID string    `json:"customer_id"`
	Amount     float64   `json:"amount"`
	CreatedAt  time.Time `json:"created_at"`
}

func (p OrderPlaced) Name() string {
	return "OrderPlaced"
}
