package app

import (
	"ecommerce/Auth/routes"
	api "ecommerce/Order/API"
	"ecommerce/Order/Domain/services"
	messaging "ecommerce/Order/Internal/Adapters"
	config "ecommerce/SharedKernel"

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

	orderPublisher := messaging.NewNatsOrderAdapter(ai.natsConn)
	orderService := services.NewOrderService(orderPublisher)
	orderModule := api.NewOrderHandler(orderService)

	authModule := routes.NewAuthHandler()

	app.RegisterModule(authModule)
	app.RegisterModule(orderModule)

	return app, nil
}
