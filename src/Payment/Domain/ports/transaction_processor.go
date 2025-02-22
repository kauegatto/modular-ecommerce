package ports

import (
	"context"
	"ecommerce/Payment/Domain/models"
)

type TransactionProcessor interface {
	Capture(ctx context.Context, card *models.Card, payment *models.Payment) (CaptureTransactionResponse, error)
	RequestCancellation(ctx context.Context, externalPaymentId string) error
}

type CaptureTransactionResponse struct {
	ExternalTransactionId string
}
