package ports

import "ecommerce/Order/Domain/models"

type OrderRepository interface {
	GetOrderById(id models.OrderID) (models.Order, error)
	Create(order models.Order) error
	PutItem(order models.Order) error
}
