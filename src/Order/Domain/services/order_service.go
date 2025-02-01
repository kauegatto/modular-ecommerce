package services

import (
	domain "ecommerce/Order/Domain"
	"ecommerce/Order/Domain/IntegrationEvents/incoming"
	"ecommerce/Order/Domain/IntegrationEvents/outgoing"
	"ecommerce/SharedKernel/models"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

type OrderService struct {
	eventBus models.Eventbus
	logger   *log.Logger
}

func NewOrderService(eventBus models.Eventbus, logger *log.Logger) (*OrderService, error) {
	service := &OrderService{
		eventBus: eventBus,
		logger:   logger,
	}

	if err := service.subscribeToEvents(); err != nil {
		return nil, fmt.Errorf("subscribe to events: %w", err)
	}

	return service, nil
}

func (s *OrderService) subscribeToEvents() error {
	s.eventBus.Subscribe(&incoming.PaymentCompleted{}, s.handlePaymentCompleted)
	s.eventBus.Subscribe(&outgoing.OrderCancelled{}, s.handleOrderCancelled)
	return nil
}

func (s *OrderService) handlePaymentCompleted(event models.Event) error {
	payment, ok := event.(*incoming.PaymentCompleted)
	if !ok {
		return fmt.Errorf("expected PaymentCompleted, got %T", event)
	}
	s.logger.Printf("[INFO] Processing payment completed for orderId %s", payment.OrderID)
	return nil
}

func (s *OrderService) handleOrderCancelled(event models.Event) error {
	orderCancelled, ok := event.(*outgoing.OrderCancelled)
	if !ok {
		return fmt.Errorf("expected PaymentCompleted, got %T", event)
	}

	s.logger.Printf("[TEST] Processed Order Cancelled Event - Order_id: %s", orderCancelled.OrderID)
	return nil
}

func (s *OrderService) PlaceOrder(customerID string, amount domain.Money) error {
	orderID, err := uuid.NewV7()
	if err != nil {
		return fmt.Errorf("generate order ID: %w", err)
	}

	event := &outgoing.OrderPlaced{
		OrderID:    orderID.String(),
		CustomerID: customerID,
		Amount:     float64(amount),
		CreatedAt:  time.Now(),
	}

	if err := s.eventBus.Publish(event); err != nil {
		return fmt.Errorf("publish order placed event: %w", err)
	}

	cancelledEvent := &outgoing.OrderCancelled{
		OrderID:     orderID.String(),
		CancelledAt: time.Now(),
		Reason:      "Test",
	}

	if err := s.eventBus.Publish(cancelledEvent); err != nil {
		return fmt.Errorf("cancelled order placed event: %w", err)
	}
	return nil
}

func (s *OrderService) CancelOrder(orderID domain.OrderID, reasib string) error {
	cancelledEvent := &outgoing.OrderCancelled{
		OrderID:     orderID.String(),
		CancelledAt: time.Now(),
		Reason:      "Test",
	}

	if err := s.eventBus.Publish(cancelledEvent); err != nil {
		return fmt.Errorf("cancelled order placed event: %w", err)
	}
	return nil
}
