package models

type Event interface {
	Name() string
}

type Eventbus interface {
	Publish(event Event) error
	Subscribe(event Event, handler func(event Event) error) error
}
