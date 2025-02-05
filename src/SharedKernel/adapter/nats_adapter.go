package adapter

import (
	"ecommerce/SharedKernel/eventBus"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/nats-io/nats.go"
)

type NatsEventbusAdapter struct {
	nc       *nats.Conn
	handlers map[string][]func(eventBus.Event) error
	mu       sync.RWMutex
}

func NewNatsEventbusAdapter(nc *nats.Conn) *NatsEventbusAdapter {
	return &NatsEventbusAdapter{
		nc:       nc,
		handlers: make(map[string][]func(eventBus.Event) error),
	}
}

func (n *NatsEventbusAdapter) Publish(event eventBus.Event) error {
	log := log.Default()
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Unable to send message of name %s. Serialization failed: %v", event.Name(), err)
		return err
	}

	err = n.nc.Publish(event.Name(), data)
	if err != nil {
		log.Printf("Unable to send message of name %s. Publish failed: %v", event.Name(), err)
		return err
	}
	return nil
}

func (n *NatsEventbusAdapter) Subscribe(event eventBus.Event, handler func(event eventBus.Event) error) error {

	n.subscribeOnNats(event, handler)
	return nil
}

func (n *NatsEventbusAdapter) subscribeOnNats(event eventBus.Event, handler func(event eventBus.Event) error) error {
	n.mu.Lock()
	n.handlers[event.Name()] = append(n.handlers[event.Name()], handler)
	n.mu.Unlock()

	if len(n.handlers[event.Name()]) == 1 {
		subscription, err := n.nc.Subscribe(event.Name(), func(msg *nats.Msg) {
			newEvent := event
			if err := json.Unmarshal(msg.Data, newEvent); err != nil {
				log.Printf("Failed to unmarshal event: %v", err)
				return
			}

			n.mu.RLock()
			handlers := n.handlers[event.Name()]
			n.mu.RUnlock()

			for _, h := range handlers {
				go func(handlerFunc func(eventBus.Event) error) {
					if err := handlerFunc(newEvent); err != nil {
						log.Printf("Handler failed for event %s: %v", event.Name(), err)
					}
				}(h)
			}
		})

		if err != nil {
			return err
		}

		if !subscription.IsValid() {
			return fmt.Errorf("invalid subscription")
		}
	}

	return nil
}
