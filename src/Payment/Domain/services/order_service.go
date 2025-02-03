package services

import (
	"ecommerce/Payment/Domain/IntegrationEvents/incoming"
	"ecommerce/SharedKernel/models"
	"fmt"
	"log"
)

type PaymentService struct {
	eventBus models.Eventbus
	logger   *log.Logger
}

func NewPaymentService(eventBus models.Eventbus, logger *log.Logger) (*PaymentService, error) {
	service := &PaymentService{
		eventBus: eventBus,
		logger:   logger,
	}

	if err := service.subscribeToEvents(); err != nil {
		return nil, fmt.Errorf("subscribe to events: %w", err)
	}

	return service, nil
}

func (s *PaymentService) subscribeToEvents() error {
	s.eventBus.Subscribe(&incoming.OrderCancelled{}, s.handleOrderCancelled)
	s.eventBus.Subscribe(&incoming.OrderPlaced{}, s.handleOrderPlaced)
	return nil
}

func (s *PaymentService) handleOrderCancelled(event models.Event) error {
	// refund payments if paid
	// disallow payment of the order

	payment, ok := event.(*incoming.OrderCancelled)
	if !ok {
		return fmt.Errorf("expected OrderCancelled, got %T", event)
	}
	s.logger.Printf("[INFO] Order Cancelled for PaymentId %s", payment.OrderID)
	return nil
}

func (s *PaymentService) handleOrderPlaced(event models.Event) error {
	OrderPlaced, ok := event.(*incoming.OrderPlaced)
	if !ok {
		return fmt.Errorf("expected OrderPlaced, got %T", event)
	}

	s.logger.Printf("[TEST] Processed OrderPlaced Event - Order: %s", OrderPlaced.OrderID)

	// create payment

	// send payment_requested

	return nil
}
