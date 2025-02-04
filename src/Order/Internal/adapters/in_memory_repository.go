package adapters

import (
	"ecommerce/Order/Domain/models"
	"fmt"
	"sync"
)

type InMemoryOrderRepository struct {
	orders map[models.OrderID]models.Order
	mutex  sync.Mutex
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{
		orders: make(map[models.OrderID]models.Order),
	}
}

func (r *InMemoryOrderRepository) GetOrderById(id models.OrderID) (models.Order, error) {
	order, ok := r.orders[id]
	if !ok {
		return models.Order{}, fmt.Errorf("order with ID %s not found", id)
	}

	return order, nil
}

func (r *InMemoryOrderRepository) Create(order models.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.orders[order.ID] = order
	return nil
}

func (r *InMemoryOrderRepository) PutItem(order models.Order) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.orders[order.ID] = order
	return nil
}
