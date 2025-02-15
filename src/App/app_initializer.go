package app

import (
	"context"
	application "ecommerce/Order/Application"
	adapters "ecommerce/Order/Infrastructure/adapters/store"
	config "ecommerce/SharedKernel"
	"ecommerce/SharedKernel/adapter"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

type AppInitializer struct {
	dbPool   *pgxpool.Pool
	config   *config.Configuration
	natsConn *nats.Conn
}

func NewAppInitializer() (*AppInitializer, error) {
	if err := config.LoadConfig(); err != nil {
		return nil, err
	}

	nc, _ := nats.Connect(config.C.NatsConfig.URL)
	dbPool, err := pgxpool.New(context.Background(), config.C.DatabaseConfig.ConnectionString())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to postgres with string %s", config.C.DatabaseConfig.ConnectionString())
	}

	return &AppInitializer{config: config.C, natsConn: nc, dbPool: dbPool}, nil
}

func (ai *AppInitializer) InitializeApp() (*App, error) {
	app := NewApp(ai.config)

	logger := log.New(os.Stdout, "[APP] ", log.LstdFlags)

	orderModule := getOrderModule(ai, logger)

	app.RegisterModule(orderModule)

	return app, nil
}

func getOrderModule(ai *AppInitializer, logger *log.Logger) *application.OrderHandler {
	natsEventBus := adapter.NewNatsEventbusAdapter(ai.natsConn)
	postgresOrderRepository := adapters.NewOrderPostgresRepository(ai.dbPool)
	orderService, err := application.NewOrderService(natsEventBus, postgresOrderRepository, logger)
	if err != nil {
		panic("error constructing order service")
	}

	orderModule := application.NewOrderHandler(orderService)
	return orderModule
}
