package normalizer

import (
	"context"
	"fmt"
	"time"

	eventsv1 "github.com/shadow0vortex/cortexops/api/v1"
	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Normalizer provides deterministic transformations from raw telemetry into strongly typed Protobuf events.
// It remains stateless and uses explicit DI if external enrichments are required in the future.
type Normalizer struct{}

func New() *Normalizer {
	return &Normalizer{}
}

// NormalizeK8sEvent converts a Kubernetes corev1.Event into a CortexOps TelemetryEnvelope.
func (n *Normalizer) NormalizeK8sEvent(ctx context.Context, k8sEvent *corev1.Event) (*eventsv1.TelemetryEnvelope, error) {
	if k8sEvent == nil {
		return nil, fmt.Errorf("received nil k8s event")
	}

	severity := "INFO"
	if k8sEvent.Type == corev1.EventTypeWarning {
		severity = "WARNING"
	}

	// Idempotent generation: we base our UUID on the K8s Event UID if available, else a new UUID.
	eventID := string(k8sEvent.UID)
	if eventID == "" {
		eventID = uuid.New().String()
	}

	metadata := &eventsv1.K8SEventMetadata{
		Namespace:       k8sEvent.InvolvedObject.Namespace,
		ResourceKind:    k8sEvent.InvolvedObject.Kind,
		ResourceName:    k8sEvent.InvolvedObject.Name,
		Uid:             string(k8sEvent.InvolvedObject.UID),
		ResourceVersion: k8sEvent.InvolvedObject.ResourceVersion,
		Action:          k8sEvent.Action,
		Reason:          k8sEvent.Reason,
		Message:         k8sEvent.Message,
	}

	envelope := &eventsv1.TelemetryEnvelope{
		EventId:   eventID,
		Timestamp: timestamppb.New(k8sEvent.LastTimestamp.Time),
		Source:    "k8s-informer",
		Severity:  severity,
		Payload: &eventsv1.TelemetryEnvelope_K8SEvent{
			K8SEvent: metadata,
		},
	}

	// Fallback timestamp if LastTimestamp was empty
	if envelope.Timestamp == nil || envelope.Timestamp.Seconds == 0 {
		envelope.Timestamp = timestamppb.New(time.Now())
	}

	return envelope, nil
}

// NormalizeMetric creates an envelope for Prometheus metrics.
func (n *Normalizer) NormalizeMetric(ctx context.Context, name string, value float64, labels map[string]string) *eventsv1.TelemetryEnvelope {
	return &eventsv1.TelemetryEnvelope{
		EventId:   uuid.New().String(),
		Timestamp: timestamppb.New(time.Now()),
		Source:    "prometheus-scraper",
		Severity:  "INFO",
		Payload: &eventsv1.TelemetryEnvelope_MetricEvent{
			MetricEvent: &eventsv1.MetricMetadata{
				MetricName: name,
				Value:      value,
				Labels:     labels,
			},
		},
	}
}
