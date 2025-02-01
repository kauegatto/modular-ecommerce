package adapter

import (
	"ecommerce/SharedKernel/models"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/nats-io/nats.go"
)

type NatsEventbusAdapter struct {
	nc     *nats.Conn
	logger *log.Logger
}

func NewNatsEventbusAdapter(nc *nats.Conn, logger *log.Logger) *NatsEventbusAdapter {
	return &NatsEventbusAdapter{
		nc:     nc,
		logger: logger,
	}
}

func (n *NatsEventbusAdapter) Publish(event models.Event) error {
	data, err := json.Marshal(event)
	if err != nil {
		n.logger.Printf("Unable to send message of name %s. Serialization failed: %v", event.Name(), err)
		return err
	}
	err = n.nc.Publish(event.Name(), data)
	if err != nil {
		n.logger.Printf("Unable to send message of name %s. Publish failed: %v", event.Name(), err)
		return err
	}
	n.logger.Printf("[EVENT_PUBLISHED] - %s", event.Name())
	return nil
}

func (n *NatsEventbusAdapter) Subscribe(event models.Event, handler func(models.Event) error) error {
	eventType := reflect.TypeOf(event)
	if eventType.Kind() == reflect.Ptr {
		eventType = eventType.Elem()
	}

	_, err := n.nc.Subscribe(event.Name(), func(msg *nats.Msg) {
		newEvent := reflect.New(eventType).Interface().(models.Event)

		if err := json.Unmarshal(msg.Data, newEvent); err != nil {
			n.logger.Printf("[ERROR] Failed to unmarshal event %s: %v", event.Name(), err)
			return
		}

		if err := handler(newEvent); err != nil {
			n.logger.Printf("[ERROR] Handler failed for event %s: %v", event.Name(), err)
			return
		}

		n.logger.Printf("[INFO] Processed event %s", event.Name())
	})

	if err != nil {
		n.logger.Printf("[ERROR] Failed to subscribe to event %s: %v", event.Name(), err)
		return fmt.Errorf("subscribe to event %s: %w", event.Name(), err)
	}

	n.logger.Printf("[INFO] Subscribed to event %s", event.Name())
	return nil
}
