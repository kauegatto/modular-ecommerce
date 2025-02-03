package domain

import (
	"time"

	"github.com/google/uuid"
)

type Money int64
type OrderID uuid.UUID

func (o OrderID) String() string {
	return uuid.UUID(o).String()
}

type Order struct {
	ID         OrderID     `json:"order_id,omitempty"`
	CustomerID string      `json:"customer_id,omitempty"`
	Status     string      `json:"status,omitempty"`
	CreatedAt  time.Time   `json:"created_at,omitempty"`
	Amount     float64     `json:"amount,omitempty"`
	OrderItems []OrderItem `json:"order_items,omitempty"`
}

type OrderItem struct {
	ProductCode string  `json:"code,omitempty"`
	Name        string  `json:"name,omitempty"`
	UnitPrice   float64 `json:"unit_price,omitempty"`
	Quantity    int     `json:"quantity,omitempty"`
}
