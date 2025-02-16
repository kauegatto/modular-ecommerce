package application

import (
	"context"
	"ecommerce/Payment/Application/IntegrationEvents/incoming"
	"ecommerce/Payment/Domain/models"
	"ecommerce/Payment/Domain/ports"
	"ecommerce/SharedKernel/eventBus"
	"fmt"
)

type PaymentService struct {
	eventBus          eventBus.Eventbus
	paymentRepository ports.PaymentRepository
}

func NewPaymentService(eventBus eventBus.Eventbus, paymentRepository ports.PaymentRepository) (*PaymentService, error) {
	service := &PaymentService{
		eventBus:          eventBus,
		paymentRepository: paymentRepository,
	}

	if err := service.subscribeToEvents(); err != nil {
		return nil, fmt.Errorf("subscribe to events: %w", err)
	}

	return service, nil
}

func (s *PaymentService) subscribeToEvents() error {
	s.eventBus.Subscribe(&incoming.OrderPlaced{}, s.handleCreatePayment)
	return nil
}

func (s *PaymentService) handleCreatePayment(event eventBus.Event) error {
	_, ok := event.(*incoming.OrderPlaced)
	if !ok {
		return fmt.Errorf("expected PaymentCompleted, got %T", event)
	}
	// todo
	return fmt.Errorf("not Implemented")
}

func (s *PaymentService) CreatePayment(ctx context.Context) (*models.Payment, error) {
	return &models.Payment{}, fmt.Errorf("not implemented")
}

func (s *PaymentService) CancelPayment(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}
func (s *PaymentService) GetPaymentById(ctx context.Context, PaymentID models.PaymentID) (*models.Payment, error) {
	return &models.Payment{}, fmt.Errorf("not implemented")
}
