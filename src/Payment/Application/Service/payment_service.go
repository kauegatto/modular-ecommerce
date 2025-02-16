package service

import (
	"context"
	"ecommerce/Payment/Application/IntegrationEvents/incoming"
	"ecommerce/Payment/Application/IntegrationEvents/outgoing"
	"ecommerce/Payment/Domain/models"
	"ecommerce/Payment/Domain/ports"
	"ecommerce/SharedKernel/eventBus"
	"fmt"
	"time"
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
	parsedEvent, ok := event.(*incoming.OrderPlaced)
	if !ok {
		return fmt.Errorf("expected PaymentCompleted, got %T", event)
	}
	payment, err := models.NewPayment(parsedEvent.OrderID, models.Money(parsedEvent.Amount))
	if err != nil {
		return fmt.Errorf("handleCreatePayment: error creating payment. Err: %v", err)
	}

	s.paymentRepository.Create(context.Background(), payment)
	return nil
}

func (s *PaymentService) ConfirmPayment(ctx context.Context, PaymentID models.PaymentID) error {
	payment, err := s.GetPaymentById(ctx, PaymentID)
	if err != nil {
		return fmt.Errorf("error getting payment %v", err)
	}
	err = payment.CompletePayment()
	if err != nil {
		return fmt.Errorf("error completing payment %v", err)
	}
	s.paymentRepository.Update(ctx, payment)

	paymentCompleted := outgoing.PaymentCompleted{
		OrderID:   payment.OrderId,
		PaymentID: PaymentID,
		Amount:    string(payment.TotalPrice()),
		Time:      time.Now(),
	}
	s.eventBus.Publish(paymentCompleted)
	return nil
}

func (s *PaymentService) CancelPayment(ctx context.Context, PaymentID models.PaymentID) error {
	payment, err := s.GetPaymentById(ctx, PaymentID)
	if err != nil {
		return fmt.Errorf("error getting payment %v", err)
	}
	err = payment.CancelPayment()
	if err != nil {
		return fmt.Errorf("error cancelling payment %v", err)
	}
	s.paymentRepository.Update(ctx, payment)

	paymentCancelled := outgoing.PaymentCancelled{
		OrderID:   payment.OrderId,
		PaymentID: PaymentID,
		Amount:    string(payment.TotalPrice()),
		Time:      time.Now(),
	}
	s.eventBus.Publish(paymentCancelled)
	return nil
}

func (s *PaymentService) GetPaymentById(ctx context.Context, PaymentID models.PaymentID) (*models.Payment, error) {
	return &models.Payment{}, fmt.Errorf("not implemented")
}
