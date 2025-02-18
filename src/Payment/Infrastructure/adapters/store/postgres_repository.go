package adapters

import (
	"context"
	"ecommerce/Payment/Domain/models"
	"ecommerce/Payment/Infrastructure/store"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db      *pgxpool.Pool
	queries *store.Queries
}

func NewPaymentPostgresRepository(db *pgxpool.Pool) PostgresRepository {
	return PostgresRepository{
		db:      db,
		queries: store.New(db),
	}
}

func (repo PostgresRepository) Create(ctx context.Context, Payment *models.Payment) error {
	request := store.CreatePaymentParams{
		ID:                   Payment.ID,
		Orderid:              Payment.OrderId,
		Totalamount:          int64(Payment.TotalPrice),
		CreatedAt:            pgtype.Timestamp{Time: Payment.CreatedAt, Valid: true},
		Integratorexternalid: pgtype.Text{String: Payment.ExternalIntegratorID, Valid: true},
	}
	_, err := repo.queries.CreatePayment(ctx, request)
	if err != nil {
		return fmt.Errorf("error creating payment %v", err)

	}
	return nil
}

func (repo PostgresRepository) Update(ctx context.Context, Payment *models.Payment) error {
	request := store.UpdatePaymentParams{
		ID:                   Payment.ID,
		Orderid:              Payment.OrderId,
		Totalamount:          int64(Payment.TotalPrice),
		CreatedAt:            pgtype.Timestamp{Time: Payment.CreatedAt},
		Integratorexternalid: pgtype.Text{String: Payment.ExternalIntegratorID},
	}
	err := repo.queries.UpdatePayment(ctx, request)
	if err != nil {
		return fmt.Errorf("error creating payment %v", err)

	}
	return nil
}

func (repo PostgresRepository) GetPaymentById(ctx context.Context, id models.PaymentID) (*models.Payment, error) {

	payment, err := repo.queries.GetPayment(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error getting payment from db %v", err)

	}

	modelPayment := paymentDbModelToModel(payment)
	return modelPayment, nil
}

func (repo PostgresRepository) GetPaymentByOrderId(ctx context.Context, orderId string) (*models.Payment, error) {

	payment, err := repo.queries.GetPaymentByOrderId(ctx, orderId)
	if err != nil {
		return nil, fmt.Errorf("error getting payment from db %v", err)
	}

	modelPayment := paymentDbModelToModel(payment)
	return modelPayment, nil
}

func paymentDbModelToModel(payment store.Payment) *models.Payment {
	modelPayment := models.NewPaymentFromRehidration(
		payment.ID,
		payment.Orderid,
		payment.Integratorexternalid.String,
		models.PaymentStatusPending,
		payment.CreatedAt.Time,
		models.Money(payment.Totalamount),
	)
	return modelPayment
}
