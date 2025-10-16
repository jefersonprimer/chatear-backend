package domain

import "context"

// Event represents an event to be published.
type Event struct {
	Subject string
	Data    []byte
}

// EventBus defines the interface for an event bus.
type EventBus interface {
	Publish(ctx context.Context, event *Event) error
}
