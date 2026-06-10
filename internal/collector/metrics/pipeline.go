package metrics

import (
	"context"
	"log/slog"

	"github.com/shadow0vortex/cortexops/internal/collector/normalizer"
	"github.com/shadow0vortex/cortexops/pkg/core"
	"google.golang.org/protobuf/proto"
)

// Pipeline manages the ingestion, normalization, and publication of Prometheus metrics.
type Pipeline struct {
	publisher  core.Publisher
	metrics    core.MetricsRecorder
	normalizer *normalizer.Normalizer
	logger     *slog.Logger
}

func NewPipeline(pub core.Publisher, metrics core.MetricsRecorder, logger *slog.Logger) *Pipeline {
	return &Pipeline{
		publisher:  pub,
		metrics:    metrics,
		normalizer: normalizer.New(),
		logger:     logger,
	}
}

// IngestRaw processes a single scraped metric and publishes it.
// In a full implementation, this would receive batches of prometheus text format data.
func (p *Pipeline) IngestRaw(ctx context.Context, metricName string, value float64, labels map[string]string) {
	p.metrics.IncCounter(ctx, "cortexops_telemetry_ingested_total", map[string]string{"source": "metrics"})

	envelope := p.normalizer.NormalizeMetric(ctx, metricName, value, labels)
	
	payload, err := proto.Marshal(envelope)
	if err != nil {
		p.logger.Error("Failed to marshal metric envelope", "error", err)
		p.metrics.IncCounter(ctx, "cortexops_telemetry_dropped_total", map[string]string{"reason": "marshal_failure"})
		return
	}

	subject := "cortex.telemetry.metrics.normalized"
	err = p.publisher.Publish(ctx, subject, envelope.EventId, payload)
	if err != nil {
		p.logger.Error("Failed to publish metric telemetry", "error", err)
		p.metrics.IncCounter(ctx, "cortexops_telemetry_dropped_total", map[string]string{"reason": "publish_failure"})
		return
	}

	p.metrics.IncCounter(ctx, "cortexops_telemetry_published_total", map[string]string{"subject": subject})
}
