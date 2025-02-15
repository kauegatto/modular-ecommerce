package ports

import (
	"context"
	"ecommerce/Order/Domain/models"
)

type OrderRepository interface {
	GetOrderById(ctx context.Context, id models.OrderID) (*models.Order, error)
	Create(ctx context.Context, order *models.Order) error
	Update(ctx context.Context, order *models.Order) error
}
