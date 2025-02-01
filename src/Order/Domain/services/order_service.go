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
	events := []models.Event{
		&incoming.PaymentCompleted{},
		&outgoing.OrderCancelled{},
	}

	for _, event := range events {
		if err := s.eventBus.Subscribe(event, s.handleEvent); err != nil {
			return fmt.Errorf("subscribe to %s: %w", event.Name(), err)
		}
	}

	return nil
}

func (s *OrderService) handleEvent(event models.Event) error {
	switch e := event.(type) {
	case *incoming.PaymentCompleted:
		return s.handlePaymentCompleted(e)
	case *outgoing.OrderCancelled:
		return s.handleOrderCancelled(e)
	default:
		return fmt.Errorf("unknown event type: %T", event)
	}
}

func (s *OrderService) handlePaymentCompleted(event *incoming.PaymentCompleted) error {
	s.logger.Printf("[INFO] Processing payment completed for order %s", event.OrderID)
	return nil
}

func (s *OrderService) handleOrderCancelled(event *outgoing.OrderCancelled) error {
	s.logger.Printf("[TEST] Processed Order Cancelled Event - Order_id: %s", event.OrderID)
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
