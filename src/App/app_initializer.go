package app

import (
	"ecommerce/Auth/routes"
	api "ecommerce/Order/API"
	"ecommerce/Order/Domain/services"
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

	natsEventBus := adapter.NewNatsEventbusAdapter(ai.natsConn)

	orderService, err := services.NewOrderService(natsEventBus, logger)
	if err != nil {
		panic("error constructing order service")
	}

	orderModule := api.NewOrderHandler(orderService)

	authModule := routes.NewAuthHandler()

	app.RegisterModule(authModule)
	app.RegisterModule(orderModule)

	return app, nil
}
