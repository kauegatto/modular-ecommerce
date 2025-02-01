package main

import (
	app "ecommerce/App"
	"log"
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

	if err != nil {
		log.Fatalf("failed to connect to NATS: %v", err)
	}

	if err := app.Start(); err != nil {
		log.Fatal("Failed to start app:", err)
	}
}
