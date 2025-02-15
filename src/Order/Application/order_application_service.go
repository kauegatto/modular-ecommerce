package application

import (
	"context"
	"ecommerce/Order/Application/IntegrationEvents/incoming"
	"ecommerce/Order/Application/IntegrationEvents/outgoing"
	"ecommerce/Order/Domain/models"
	"ecommerce/Order/Domain/ports"
	"ecommerce/SharedKernel/eventBus"
	"fmt"
	"log"
	"time"
)

type OrderService struct {
	eventBus        eventBus.Eventbus
	orderRepository ports.OrderRepository
	logger          *log.Logger
}

func NewOrderService(eventBus eventBus.Eventbus, orderRepository ports.OrderRepository, logger *log.Logger) (*OrderService, error) {
	service := &OrderService{
		eventBus:        eventBus,
		orderRepository: orderRepository,
		logger:          logger,
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

func (s *OrderService) handlePaymentCompleted(event eventBus.Event) error {
	payment, ok := event.(*incoming.PaymentCompleted)
	if !ok {
		return fmt.Errorf("expected PaymentCompleted, got %T", event)
	}
	s.logger.Printf("[INFO] Processing payment completed for orderId %s", payment.OrderID)
	return nil
}

func (s *OrderService) handleOrderCancelled(event eventBus.Event) error {
	orderCancelled, ok := event.(*outgoing.OrderCancelled)
	if !ok {
		return fmt.Errorf("expected PaymentCompleted, got %T", event)
	}

	s.logger.Printf("[TEST] Processed Order Cancelled Event - Order_id: %s", orderCancelled.OrderID)
	return nil
}

func (s *OrderService) PlaceOrder(ctx context.Context, customerID string, items []models.OrderItem) (models.Order, error) {
	order, err := models.NewOrder(customerID, items)
	if err != nil {
		return models.Order{}, fmt.Errorf("generate order ID: %w", err)
	}

	if err := s.orderRepository.Create(ctx, order); err != nil {
		return models.Order{}, fmt.Errorf("order creation failed %w", err)
	}

	event := &outgoing.OrderPlaced{
		OrderID:    order.ID.String(),
		CustomerID: customerID,
		Amount:     float64(order.TotalPrice()),
		CreatedAt:  time.Now(),
	}

	if err := s.eventBus.Publish(event); err != nil {
		return models.Order{}, fmt.Errorf("publish order placed event: %w", err)
	}

	return *order, nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID models.OrderID, reason string) error {
	order, err := s.GetOrderById(ctx, orderID)
	if err != nil {
		return fmt.Errorf("cancelled order placed event: %w", err)
	}
	order.CancelOrder()
	s.orderRepository.Update(ctx, &order)

	cancelledEvent := &outgoing.OrderCancelled{
		OrderID:     orderID.String(),
		CancelledAt: time.Now(),
		Reason:      reason,
	}

	if err := s.eventBus.Publish(cancelledEvent); err != nil {
		return fmt.Errorf("cancelled order placed event: %w", err)
	}
	return nil
}
func (s *OrderService) GetOrderById(ctx context.Context, orderID models.OrderID) (models.Order, error) {
	order, err := s.orderRepository.GetOrderById(ctx, orderID)
	if err != nil {
		return models.Order{}, fmt.Errorf("error getting order: %w", err)
	}
	return *order, nil
}
