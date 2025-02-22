package adapters

import (
	"context"
	"ecommerce/Payment/Domain/models"
	"ecommerce/Payment/Infrastructure/store"
	"fmt"
	"strings"

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
	kind, err := repo.getDbKind(ctx, &Payment.Kind)
	if err != nil {
		return fmt.Errorf("error getting paymentKind from db %v", err)
	}

	status, err := repo.getDbStatus(ctx, &Payment.Status)
	if err != nil {
		return fmt.Errorf("error getting paymentStatus from db %v", err)
	}

	request := store.CreatePaymentParams{
		ID:                   Payment.ID,
		Orderid:              Payment.OrderId,
		Totalamount:          int64(Payment.TotalPrice),
		CreatedAt:            pgtype.Timestamp{Time: Payment.CreatedAt, Valid: true},
		Integratorexternalid: pgtype.Text{String: Payment.ExternalIntegratorID, Valid: true},
		KindID:               kind.ID,
		StatusID:             status.ID,
	}

	_, err = repo.queries.CreatePayment(ctx, request)
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
		KindID:               0,
		StatusID:             0,
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

	modelPayment, err := repo.paymentDbModelToModel(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("error reconstructing payment from db %v", err)

	}
	return modelPayment, nil
}

func (repo PostgresRepository) GetPaymentByOrderId(ctx context.Context, orderId string) (*models.Payment, error) {

	payment, err := repo.queries.GetPaymentByOrderId(ctx, orderId)

	if err != nil {
		return nil, fmt.Errorf("error getting payment from db %v", err)
	}

	modelPayment, err := repo.paymentDbModelToModel(ctx, payment)
	if err != nil {
		return nil, fmt.Errorf("error reconstructing payment from db %v", err)

	}
	return modelPayment, nil
}

func (repo PostgresRepository) getPaymentKind(ctx context.Context, kindId int32) (*models.PaymentKind, error) {
	kind, err := repo.queries.GetKindById(ctx, kindId)
	if err != nil {
		return nil, fmt.Errorf("error getting payment from db %v", err)
	}
	modelKind := models.PaymentKind(kind.Name) // i wish i had an error if it was not valid
	return &modelKind, nil
}

func (repo PostgresRepository) getPaymentStatus(ctx context.Context, statusID int32) (*models.PaymentStatus, error) {
	status, err := repo.queries.GetStatusById(ctx, statusID)
	if err != nil {
		return nil, fmt.Errorf("error getting payment from db %v", err)
	}
	modelstatus := models.PaymentStatus(status.Name) // i wish i had an error if it was not valid
	return &modelstatus, nil
}

func (repo PostgresRepository) getDbKind(ctx context.Context, kind *models.PaymentKind) (*store.PaymentKind, error) {
	allKinds, err := repo.queries.GetPaymentKind(ctx) // estava com preguiça de fazer isso melhor, todo
	if err != nil {
		return nil, fmt.Errorf("error getting payment kinds from db: %w", err)
	}

	for _, k := range allKinds {
		if strings.EqualFold(k.Name, string(*kind)) {
			return &k, nil
		}
	}

	return nil, fmt.Errorf("kind not found: %s", *kind)
}

func (repo PostgresRepository) getDbStatus(ctx context.Context, status *models.PaymentStatus) (*store.PaymentStatus, error) {
	allStatuses, err := repo.queries.ListStatus(ctx) // estava com preguiça de fazer isso melhor, todo
	if err != nil {
		return nil, fmt.Errorf("error getting payment statuses from db: %w", err)
	}

	for _, s := range allStatuses {
		if strings.EqualFold(s.Name, string(*status)) {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("status not found: %s", *status)
}
func (repo PostgresRepository) paymentDbModelToModel(ctx context.Context, payment store.Payment) (*models.Payment, error) {
	kind, err := repo.getPaymentKind(ctx, payment.KindID)
	if err != nil {
		return nil, fmt.Errorf("error getting paymentKind from db %v", err)
	}

	status, err := repo.getPaymentStatus(ctx, payment.StatusID)
	if err != nil {
		return nil, fmt.Errorf("error getting paymentStatus from db %v", err)
	}

	modelPayment := models.NewPaymentFromRehidration(
		payment.ID,
		payment.Orderid,
		payment.Integratorexternalid.String,
		*status,
		payment.CreatedAt.Time,
		models.Money(payment.Totalamount),
		*kind,
	)
	return modelPayment, nil
}
