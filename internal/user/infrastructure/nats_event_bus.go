package infrastructure

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
	"github.com/nats-io/nats.go"
)

const (
	maxRetries    = 3
	retryInterval = 1 * time.Second
)

// NATSEventBus is a NATS implementation of the domain.EventBus.
type NATSEventBus struct {
	Conn *nats.Conn
}

// NewNATSEventBus creates a new NATSEventBus.
func NewNATSEventBus(conn *nats.Conn) *NATSEventBus {
	return &NATSEventBus{Conn: conn}
}

// Publish publishes an event to the event bus with retry mechanism.
func (b *NATSEventBus) Publish(ctx context.Context, event *domain.Event) error {
	for i := 0; i < maxRetries; i++ {
		err := b.Conn.Publish(event.Subject, event.Data)
		if err == nil {
			return nil
		}
		log.Printf("Attempt %d to publish event to subject %s failed: %v", i+1, event.Subject, err)
		if i < maxRetries-1 {
			time.Sleep(retryInterval)
		}
	}
	return fmt.Errorf("failed to publish event to subject %s after %d attempts", event.Subject, maxRetries)
}
