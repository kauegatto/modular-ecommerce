package app

import (
	application "ecommerce/Order/Application"
	"ecommerce/Order/Internal/adapters"
	config "ecommerce/SharedKernel"
	"ecommerce/SharedKernel/adapter"
	"log"
	"os"

	"github.com/nats-io/nats.go"
)

type AppInitializer struct {
	config   *config.NatsConfig
	natsConn *nats.Conn
}

func NewAppInitializer() (*AppInitializer, error) {
	if err := config.LoadConfig(); err != nil {
		return nil, err
	}

	nc, _ := nats.Connect(config.NatsConfiguration.URL)

	return &AppInitializer{config: config.NatsConfiguration, natsConn: nc}, nil
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
	inMemOrderRepository := adapters.NewInMemoryOrderRepository()
	orderService, err := application.NewOrderService(natsEventBus, inMemOrderRepository, logger)
	if err != nil {
		panic("error constructing order service")
	}

	orderModule := application.NewOrderHandler(orderService)
	return orderModule
}
