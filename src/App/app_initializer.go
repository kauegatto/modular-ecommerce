package app

import (
	"context"
	orderapplication "ecommerce/Order/Application"
	orderadapters "ecommerce/Order/Infrastructure/adapters/store"
	application "ecommerce/Payment/Application/Routing"
	paymentapplication "ecommerce/Payment/Application/Routing"
	paymentservice "ecommerce/Payment/Application/Service"
	paymentadapters "ecommerce/Payment/Infrastructure/adapters"

	config "ecommerce/SharedKernel"
	"ecommerce/SharedKernel/adapter"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

type AppInitializer struct {
	dbPool       *pgxpool.Pool
	config       *config.Configuration
	NatsEventBus *adapter.NatsEventbusAdapter
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

	eventBus := adapter.NewNatsEventbusAdapter(nc)
	return &AppInitializer{config: config.C, dbPool: dbPool, NatsEventBus: eventBus}, nil
}

func (ai *AppInitializer) InitializeApp() (*App, error) {
	app := NewApp(ai.config)

	orderModule := ai.getOrderModule()
	paymentModule := ai.getPaymentModule()

	app.RegisterModule(orderModule)
	app.RegisterModule(paymentModule)

	return app, nil
}

func (ai *AppInitializer) getOrderModule() *orderapplication.OrderHandler {
	postgresOrderRepository := orderadapters.NewOrderPostgresRepository(ai.dbPool)
	orderService, err := orderapplication.NewOrderService(ai.NatsEventBus, postgresOrderRepository)
	if err != nil {
		panic("error constructing order service")
	}

	orderModule := orderapplication.NewOrderHandler(orderService)
	return orderModule
}

func (ai *AppInitializer) getPaymentModule() *paymentapplication.PaymentHandler {
	eRedeConfig := paymentadapters.ERedeConfig{
		PV:      "",
		Token:   "",
		BaseURL: "",
		Timeout: 0,
	}

	postgresOrderRepository := paymentadapters.NewPaymentPostgresRepository(ai.dbPool)
	eRedePaymentProcessor := paymentadapters.NewERedeProcessor(eRedeConfig)
	paymentService, err := paymentservice.NewPaymentService(ai.NatsEventBus, postgresOrderRepository, eRedePaymentProcessor)
	if err != nil {
		panic("error constructing order service")
	}
	return application.NewPaymentHandler(paymentService)
}
