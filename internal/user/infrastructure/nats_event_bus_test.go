package infrastructure

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/jefersonprimer/chatear-backend/internal/user/domain"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockNATSConn is a mock implementation of *nats.Conn for testing
type MockNATSConn struct {
	mock.Mock
}

// Publish implements the Publish method of *nats.Conn
func (m *MockNATSConn) Publish(subj string, data []byte) error {
	args := m.Called(subj, data)
	return args.Error(0)
}

func TestNATSEventBus_Publish(t *testing.T) {
	ctx := context.Background()
	event := &domain.Event{
		Subject: "test.subject",
		Data:    []byte("test data"),
	}

	// Test case 1: Successful publish
	mockNATS := new(MockNATSConn)
	mockNATS.On("Publish", event.Subject, event.Data).Return(nil).Once()
	bus := NewNATSEventBus(mockNATS)
	err := bus.Publish(ctx, event)
	assert.NoError(t, err)
	mockNATS.AssertExpectations(t)

	// Test case 2: Publish fails once, then succeeds on retry
	mockNATS = new(MockNATSConn)
	mockNATS.On("Publish", event.Subject, event.Data).Return(errors.New("nats error")).Once()
	mockNATS.On("Publish", event.Subject, event.Data).Return(nil).Once()
	bus = NewNATSEventBus(mockNATS)
	err = bus.Publish(ctx, event)
	assert.NoError(t, err)
	mockNATS.AssertExpectations(t)

	// Test case 3: Publish fails multiple times
	mockNATS = new(MockNATSConn)
	mockNATS.On("Publish", event.Subject, event.Data).Return(errors.New("nats error")).Times(maxRetries)
	bus = NewNATSEventBus(mockNATS)
	err = bus.Publish(ctx, event)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to publish event")
	mockNATS.AssertExpectations(t)
}
