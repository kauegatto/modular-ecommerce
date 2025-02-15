package adapters

import (
	"context"
	"ecommerce/Order/Domain/models"
	"ecommerce/Order/Infrastructure/store"
	"math/big"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db      *pgxpool.Pool
	queries *store.Queries
}

func NewOrderPostgresRepository(db *pgxpool.Pool) PostgresRepository {
	return PostgresRepository{db: db, queries: store.New(db)}
}

func (repo PostgresRepository) Create(ctx context.Context, order *models.Order) error {
	customerUUID, err := uuid.Parse(order.CustomerID)
	if err != nil {
		return err
	}

	params := store.CreateOrderParams{
		ID:         order.ID,
		CustomerID: customerUUID,
		StatusID:   1,
		CreatedAt:  pgtype.Timestamp{Time: time.Now(), Valid: true},
		TotalPrice: int64(order.TotalPrice()),
		Discount:   pgtype.Numeric{Int: big.NewInt(int64(order.Discount())), Valid: true},
	}

	_, err = repo.queries.CreateOrder(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (repo PostgresRepository) Update(ctx context.Context, order *models.Order) error {
	params := store.UpdateOrderParams{
		ID:         order.ID,
		StatusID:   1,
		TotalPrice: int64(order.TotalPrice()),
		Discount:   pgtype.Numeric{Int: big.NewInt(int64(order.Discount())), Valid: true},
	}

	err := repo.queries.UpdateOrder(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func (repo PostgresRepository) GetOrderById(ctx context.Context, id models.OrderID) (*models.Order, error) {
	dbOrder, err := repo.queries.GetOrder(ctx, id)
	if err != nil {
		return &models.Order{}, err
	}
	order := models.NewOrderFromDTO(
		dbOrder.ID,
		dbOrder.CustomerID.String(),
		models.OrderStatus(dbOrder.StatusID),
		dbOrder.CreatedAt.Time,
		float64(dbOrder.TotalPrice),
		int(dbOrder.Discount.Int.Int64()), // todo
		[]models.OrderItem{},              // todo
	)

	return order, nil
}
