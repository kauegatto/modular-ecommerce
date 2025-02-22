package adapters

import (
	"context"
	"ecommerce/Order/Domain/models"
	"ecommerce/Order/Infrastructure/store"
	"fmt"
	"log/slog"
	"math/big"
	"strings"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	db      *pgxpool.Pool
	queries *store.Queries
}

func NewOrderPostgresRepository(db *pgxpool.Pool) PostgresRepository {
	return PostgresRepository{
		db:      db,
		queries: store.New(db),
	}
}

func (repo PostgresRepository) Create(ctx context.Context, order *models.Order) error {
	params := store.CreateOrderParams{
		ID:         order.ID,
		CustomerID: order.CustomerID,
		StatusID:   1,
		CreatedAt:  pgtype.Timestamp{Time: time.Now(), Valid: true},
		TotalPrice: int64(order.TotalPrice()),
		Discount:   pgtype.Numeric{Int: big.NewInt(int64(order.Discount())), Valid: true},
	}

	_, err := repo.queries.CreateOrder(ctx, params)
	if err != nil {
		slog.ErrorContext(ctx, "Order when inserting order on db %w", slog.Attr{Key: "error", Value: slog.AnyValue(err)})
		return err
	}
	slog.InfoContext(ctx, "Order inserted on db with success", slog.Attr{Key: "order", Value: slog.AnyValue(order)})
	return nil
}

func (repo PostgresRepository) Update(ctx context.Context, order *models.Order) error {
	status, err := repo.getDbStatus(ctx, order.Status())
	if err != nil {
		slog.ErrorContext(ctx, "Order when getting orderStatus id on db",
			slog.Attr{Key: "error", Value: slog.AnyValue(err)},
			slog.Attr{Key: "Status", Value: slog.StringValue(string(order.Status()))})
		return err
	}

	params := store.UpdateOrderParams{
		ID:         order.ID,
		StatusID:   status.ID,
		TotalPrice: int64(order.TotalPrice()),
		Discount:   pgtype.Numeric{Int: big.NewInt(int64(order.Discount())), Valid: true},
	}

	err = repo.queries.UpdateOrder(ctx, params)
	if err != nil {
		slog.ErrorContext(ctx, "Order when updating order on db", slog.Attr{Key: "error", Value: slog.AnyValue(err)})
		return err
	}
	slog.InfoContext(ctx, "Order updated on db with success %w", slog.Attr{Key: "order", Value: slog.AnyValue(order)})
	return nil
}

func (repo PostgresRepository) GetOrderById(ctx context.Context, id models.OrderID) (*models.Order, error) {
	dbOrder, err := repo.queries.GetOrder(ctx, id)
	if err != nil {
		slog.ErrorContext(ctx, "Could not find order", slog.Attr{Key: "OrderId", Value: slog.StringValue(id.String())})
		return &models.Order{}, err
	}

	status, err := repo.queries.GetStatus(ctx, dbOrder.StatusID)
	if err != nil {
		slog.ErrorContext(ctx, "Could not find status", slog.Attr{Key: "StatusId", Value: slog.AnyValue(dbOrder.StatusID)})
		return &models.Order{}, err
	}

	order := models.NewOrderFromDTO(
		dbOrder.ID,
		dbOrder.CustomerID,
		models.OrderStatus(status.StatusName),
		dbOrder.CreatedAt.Time,
		float64(dbOrder.TotalPrice),
		int(dbOrder.Discount.Int.Int64()), // todo
		[]models.OrderItem{},              // todo
	)

	slog.InfoContext(ctx, "Order found", slog.Attr{Key: "OrderId", Value: slog.StringValue(id.String())})
	return order, nil
}

func (repo PostgresRepository) getDbStatus(ctx context.Context, status models.OrderStatus) (*store.Status, error) {
	allStatuses, err := repo.queries.ListStatuses(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting payment statuses from db: %w", err)
	}

	for _, s := range allStatuses {
		if strings.EqualFold(s.StatusName, string(status)) {
			return &s, nil
		}
	}

	return nil, fmt.Errorf("status not found: %s", status)
}
