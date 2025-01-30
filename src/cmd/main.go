package main

import (
	app "ecommerce/App"
	"log"

	"github.com/nats-io/nats.go"
)

func main() {
	appInitializer, err := app.NewAppInitializer()
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}

	app, err := appInitializer.InitializeApp()
	if err != nil {
		log.Fatal("Failed to initialize app:", err)
	}

	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	if err := app.Start(); err != nil {
		log.Fatal("Failed to start app:", err)
	}
}
