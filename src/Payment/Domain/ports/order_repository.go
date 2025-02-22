package ports

import (
	"context"
	"ecommerce/Payment/Domain/models"
)

type PaymentRepository interface {
	GetPaymentById(ctx context.Context, id models.PaymentID) (*models.Payment, error)
	GetPaymentByOrderId(ctx context.Context, orderId string) (*models.Payment, error)
	Create(ctx context.Context, Payment *models.Payment) error
	Update(ctx context.Context, Payment *models.Payment) error
}
