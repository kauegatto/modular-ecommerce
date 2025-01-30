package messaging

import (
	"ecommerce/Order/Domain/events/incoming"
	"ecommerce/Order/Domain/events/outgoing"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type NatsOrderAdapter struct {
	nc *nats.Conn
}

func NewNatsOrderAdapter(nc *nats.Conn) *NatsOrderAdapter {
	return &NatsOrderAdapter{nc: nc}
}

func (n *NatsOrderAdapter) PublishOrderPlaced(event outgoing.OrderPlaced) error {
	n.nc.Publish("domain.payment.completed", []byte("[TEST] [todo remove] hello"))
	data, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal OrderPlaced event: %w", err)
	}
	return n.nc.Publish("domain.order.placed", data)
}

func (n *NatsOrderAdapter) SubscribeToPaymentCompleted(handler func(incoming.PaymentCompleted) error) error {
	_, err := n.nc.Subscribe("domain.payment.completed", func(msg *nats.Msg) {
		log.Println("[NATS] Received payment completed event:", string(msg.Data))

		var event incoming.PaymentCompleted
		if err := json.Unmarshal(msg.Data, &event); err != nil {
			log.Println("Error unmarshalling PaymentCompleted event:", err)
			return
		}

		if err := handler(event); err != nil {
			log.Println("Error processing PaymentCompleted event:", err)
		}
	})
	return err
}
