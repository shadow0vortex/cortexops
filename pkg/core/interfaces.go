package core

import (
	"context"
)

// Publisher defines the contract for emitting events to the message broker.
type Publisher interface {
	Publish(ctx context.Context, subject string, eventID string, payload []byte) error
}

// Subscriber defines the contract for consuming events from the message broker.
type Subscriber interface {
	Subscribe(ctx context.Context, subject string, handler EventHandler) error
	Close() error
}

// EventHandler is the callback invoked when a subscriber receives a message.
type EventHandler func(ctx context.Context, payload []byte) error

// EventSerializer handles deterministic serialization and deserialization of Protobuf events.
type EventSerializer interface {
	Marshal(v interface{}) ([]byte, error)
	Unmarshal(data []byte, v interface{}) error
}

// TraceProvider abstracts distributed tracing mechanisms (e.g., OpenTelemetry).
type TraceProvider interface {
	StartSpan(ctx context.Context, name string, attributes map[string]string) (context.Context, Span)
}

// Span represents a single operation within a trace.
type Span interface {
	End()
	RecordError(err error)
	SetAttributes(attributes map[string]string)
}

// MetricsRecorder defines the contract for exposing observability metrics (e.g., Prometheus).
type MetricsRecorder interface {
	IncCounter(ctx context.Context, name string, labels map[string]string)
	ObserveHistogram(ctx context.Context, name string, value float64, labels map[string]string)
	SetGauge(ctx context.Context, name string, value float64, labels map[string]string)
}

// PolicyEvaluator evaluates deterministic governance rules (e.g., via Open Policy Agent).
type PolicyEvaluator interface {
	EvaluateRemediation(ctx context.Context, incidentID string, requestedAction string) (bool, error)
}

// IncidentStore abstracts persistent storage for incident lifecycles (e.g., PostgreSQL).
type IncidentStore interface {
	SaveIncident(ctx context.Context, incidentID string, payload []byte) error
	GetIncident(ctx context.Context, incidentID string) ([]byte, error)
	UpdateState(ctx context.Context, incidentID string, newState string) error
}

// TopologyProvider abstracts the querying of service dependencies and blast radii (e.g., GraphDB).
type TopologyProvider interface {
	GetDependencies(ctx context.Context, serviceName string) ([]string, error)
	CalculateBlastRadius(ctx context.Context, serviceName string) ([]string, error)
}
