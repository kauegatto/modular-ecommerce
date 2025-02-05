package dto

import (
	"ecommerce/Order/Domain/models"
	"time"
)

type OrderDto struct {
	ID         models.OrderID     `json:"order_id,omitempty"`
	CustomerID string             `json:"customer_id,omitempty"`
	Status     models.OrderStatus `json:"status,omitempty"`
	CreatedAt  time.Time          `json:"createdAt,omitempty"`
	TotalPrice float64            `json:"totalPrice,omitempty"`
	Items      []models.OrderItem `json:"items,omitempty"`
}

func ToOrderDto(order *models.Order) OrderDto {
	return OrderDto{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		Status:     order.Status(),
		CreatedAt:  order.CreatedAt(),
		TotalPrice: order.TotalPrice(),
		Items:      order.Items(),
	}
}
