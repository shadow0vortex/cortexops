package broker

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

// Broker represents a generic interface for message bus interactions.
type Broker interface {
	Publish(ctx context.Context, subject string, data []byte) error
	Subscribe(ctx context.Context, subject string, handler Handler) error
	Close() error
}

// Handler is the callback function for processing incoming messages.
type Handler func(ctx context.Context, data []byte) error

// MockBroker is a simple implementation for testing or local development without NATS.
type MockBroker struct {
	logger *slog.Logger
}

func NewMockBroker(logger *slog.Logger) *MockBroker {
	return &MockBroker{logger: logger}
}

func (m *MockBroker) Publish(ctx context.Context, subject string, data []byte) error {
	m.logger.Info("Mock publish", "subject", subject, "data_size", len(data))
	return nil
}

func (m *MockBroker) Subscribe(ctx context.Context, subject string, handler Handler) error {
	m.logger.Info("Mock subscribe", "subject", subject)
	return nil
}

func (m *MockBroker) Close() error {
	return nil
}

// Note: A real NATS JetStream implementation would go here, utilizing "github.com/nats-io/nats.go"
// To keep Phase 1 scope manageable and ensure it compiles without a NATS server locally, 
// the interface is defined with a mock. The real implementation uses the exact same interface.
