package adapter

import (
	"ecommerce/SharedKernel/models"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
)

type NatsEventbusAdapter struct {
	nc *nats.Conn
}

func (n *NatsEventbusAdapter) Publish(event models.Event) error {
	log := log.Default()
	data, err := json.Marshal(event)
	if err != nil {
		log.Println("Unable to send message of name %s. Serialization failed %e", event.Name(), err)
	}
	err = n.nc.Publish(event.Name(), data)
	if err != nil {
		log.Println("Unable to send message of name %s. Publish failed %e", event.Name(), err)
	}
	return nil
}

func (n *NatsEventbusAdapter) Subscribe(event models.Event, handler func(event models.Event) error) error {
	subscription, err := n.nc.Subscribe(event.Name(), func(msg *nats.Msg) {
		var event models.Event
		json.Unmarshal(msg.Data, &event)
		handler(event)
	})
	if !subscription.IsValid() || err != nil {
		return err
	}
}
