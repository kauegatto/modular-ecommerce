package services

import (
	domain "ecommerce/Order/Domain"
	"ecommerce/Order/Domain/events/incoming"
	"ecommerce/Order/Domain/events/outgoing"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EventPublisher interface {
	PublishOrderPlaced(event outgoing.OrderPlaced) error
}

type OrderService struct {
	eventPublisher EventPublisher
}

func NewOrderService(publisher EventPublisher) *OrderService {
	return &OrderService{
		eventPublisher: publisher,
	}
}

func (s *OrderService) PlaceOrder(customerId string, amount domain.Money) error {
	// bussiness logic
	uuid, err := uuid.NewV7()
	if err != nil {
		return err
	}

	err = s.eventPublisher.PublishOrderPlaced(outgoing.OrderPlaced{
		OrderID:    uuid.String(),   //order.ID,
		CustomerID: customerId,      //order.CustomerID,
		Amount:     float64(amount), //order.TotalAmount, -- todo how to / where to  map from domain to shared context
		CreatedAt:  time.Now(),
	})

	if err != nil {
		return fmt.Errorf("failed to publish order placed event: %w", err)
	}

	return nil
}

func (s *OrderService) OnPaymentCompleted(incoming.PaymentCompleted) error {
	// todo bussiness logic

	err := s.eventPublisher.PublishOrderPlaced(outgoing.OrderPlaced{
		OrderID:    "id",     //order.ID,
		CustomerID: "custId", //order.CustomerID,
		Amount:     123,      //order.TotalAmount,
		CreatedAt:  time.Now(),
	})

	if err != nil {
		return fmt.Errorf("failed to publish order placed event: %w", err)
	}
	return nil
}
func (s *OrderService) HandlePaymentCompleted(event incoming.PaymentCompleted) error {
	// publish order_confirmed
	return nil
}
