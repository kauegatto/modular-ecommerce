package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Money int64
type OrderID uuid.UUID
type OrderStatus string

func (o OrderID) String() string {
	return uuid.UUID(o).String()
}

func NewOrderID(s string) (OrderID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return OrderID{}, err
	}
	return OrderID(id), nil
}

const (
	OrderStatusPlaced    OrderStatus = "ORDER_PLACED"
	OrderStatusPending   OrderStatus = "ORDER_PENDING"
	OrderStatusConfirmed OrderStatus = "ORDER_CONFIRMED"
	OrderStatusCancelled OrderStatus = "ORDER_CANCELLED"
)

type Order struct {
	ID         OrderID `json:"order_id,omitempty"`
	customerID string
	status     OrderStatus
	createdAt  time.Time
	totalPrice float64
	items      []OrderItem
}

func (order *Order) AddItems(items []OrderItem) []error {
	var errors []error
	for _, item := range items {
		if err := order.AddItem(item); err != nil {
			errors = append(errors, err)
		}
		order.items = append(order.items, item)
	}
	return errors
}

func (order *Order) AddItem(item OrderItem) error {
	order.items = append(order.items, item)
	order.totalPrice = item.UnitPrice + order.totalPrice
	return nil
}

func (o *Order) ConfirmOrder() error {
	if o.status != OrderStatusPlaced {
		return fmt.Errorf("cannot confirm order with status: %s", o.status)
	}
	o.status = OrderStatusConfirmed
	return nil
}

func (o *Order) PendingOrder() error {
	if o.status != OrderStatusPlaced {
		return fmt.Errorf("cannot set order to pending with status: %s", o.status)
	}
	o.status = OrderStatusPending
	return nil
}

func (o *Order) CancelOrder() error {
	if o.status == OrderStatusCancelled {
		return fmt.Errorf("order is already cancelled")
	}
	o.status = OrderStatusCancelled
	return nil
}

func (o *Order) TotalPrice() float64 {
	return o.totalPrice
}

func NewOrder(customerID string, items []OrderItem) (*Order, error) {
	id, _ := uuid.NewV7()
	order := &Order{
		ID:         OrderID(id),
		customerID: customerID,
		status:     OrderStatusPlaced,
		createdAt:  time.Now().UTC(),
		totalPrice: 0,
		items:      []OrderItem{},
	}

	if err := order.AddItems(items); err != nil {
		return nil, fmt.Errorf("failed to add items: %v", err)
	}

	return order, nil
}

type OrderItem struct {
	ProductCode string  `json:"code,omitempty"`
	Name        string  `json:"name,omitempty"`
	UnitPrice   float64 `json:"unit_price,omitempty"`
	Quantity    int     `json:"quantity,omitempty"`
}
