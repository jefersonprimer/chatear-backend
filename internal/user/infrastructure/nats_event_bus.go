package infrastructure

import (
	"context"

	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
	"github.com/nats-io/nats.go"
)

// NATSEventBus is a NATS implementation of the domain.EventBus.
type NATSEventBus struct {
	Conn *nats.Conn
}

// NewNATSEventBus creates a new NATSEventBus.
func NewNATSEventBus(natsURL string) (*NATSEventBus, error) {
	conn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}
	return &NATSEventBus{Conn: conn}, nil
}

// Publish publishes an event to the event bus.
func (b *NATSEventBus) Publish(ctx context.Context, event *domain.Event) error {
	return b.Conn.Publish(event.Subject, event.Data)
}
