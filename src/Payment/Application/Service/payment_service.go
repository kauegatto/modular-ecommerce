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
	eventBus             eventBus.Eventbus
	paymentRepository    ports.PaymentRepository
	transactionProcessor ports.TransactionProcessor
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

	err = s.paymentRepository.Create(context.Background(), payment)
	if err != nil {
		return fmt.Errorf("error creating payment on database %v", err)
	}

	// should have been done repository + event send transactionally using transactional outbox pattern
	paymentCreated := outgoing.PaymentCreated{
		OrderID:   payment.OrderId,
		PaymentID: payment.ID,
		Amount:    string(payment.TotalPrice),
		Time:      time.Now(),
	}
	s.eventBus.Publish(paymentCreated)
	return nil
}

func (s *PaymentService) CapturePayment(ctx context.Context, PaymentID models.PaymentID, card *models.Card) error {
	// if error on capture process, already cancel it.
	// In a real world scenario we'd just republish the event and avoid losses

	// if capture is debit and succeded (synchronous operation) -> confirm it
	// else, just wait for the callback from payment gateway
	payment, err := s.paymentRepository.GetPaymentById(ctx, PaymentID)
	if err != nil {
		return fmt.Errorf("error getting payment %v", err)
	}

	err = s.transactionProcessor.Capture(ctx, card, payment)
	if err != nil {
		err = s.CancelPayment(ctx, PaymentID)
		if err != nil {
			return err
		}
		return fmt.Errorf("error processing payment %v", err)
	}
	if payment.Kind == models.PaymentKindDebit {
		return s.ConfirmPayment(ctx, PaymentID)
	}

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

	err = s.paymentRepository.Update(ctx, payment)
	if err != nil {
		return fmt.Errorf("error completing payment on database %v", err)
	}

	paymentCompleted := outgoing.PaymentCompleted{
		OrderID:   payment.OrderId,
		PaymentID: PaymentID,
		Amount:    string(payment.TotalPrice),
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

	err = s.paymentRepository.Update(ctx, payment)
	if err != nil {
		return fmt.Errorf("error cancelling payment on database %v", err)
	}

	paymentCancelled := outgoing.PaymentCancelled{
		OrderID:   payment.OrderId,
		PaymentID: PaymentID,
		Amount:    string(payment.TotalPrice),
		Time:      time.Now(),
	}
	s.eventBus.Publish(paymentCancelled)
	return nil
}

func (s *PaymentService) GetPaymentById(ctx context.Context, PaymentID models.PaymentID) (*models.Payment, error) {
	payment, err := s.paymentRepository.GetPaymentById(ctx, PaymentID)
	if err != nil {
		return &models.Payment{}, fmt.Errorf("error finding payment on database %v", err)
	}
	return payment, nil
}

func (s *PaymentService) GetPaymentByOrderId(ctx context.Context, orderId string) (*models.Payment, error) {
	payment, err := s.paymentRepository.GetPaymentByOrderId(ctx, orderId)
	if err != nil {
		return &models.Payment{}, fmt.Errorf("error finding payment on database %v", err)
	}
	return payment, nil
}
