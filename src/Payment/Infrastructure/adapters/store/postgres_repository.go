package adapters

import (
	"context"
	"ecommerce/Payment/Domain/models"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewPaymentPostgresRepository(db *pgxpool.Pool) PostgresRepository {
	return PostgresRepository{
		db: db,
	}
}

func (repo PostgresRepository) Create(ctx context.Context, Payment *models.Payment) error {
	return fmt.Errorf("not Implemented")
}

func (repo PostgresRepository) Update(ctx context.Context, Payment *models.Payment) error {
	return fmt.Errorf("not Implemented")

}

func (repo PostgresRepository) GetPaymentById(ctx context.Context, id models.PaymentID) (*models.Payment, error) {
	return &models.Payment{}, fmt.Errorf("not Implemented")
}
