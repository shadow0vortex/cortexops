package telemetry

import (
	"context"

	"github.com/shadow0vortex/cortexops/pkg/core"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

// OTelProvider implements the core.TraceProvider interface.
type OTelProvider struct {
	tracer trace.Tracer
}

// NewOTelProvider creates a new trace provider for a specific service.
// In a full implementation, this would configure the OTLP exporter.
func NewOTelProvider(serviceName string) *OTelProvider {
	return &OTelProvider{
		tracer: otel.Tracer(serviceName),
	}
}

// StartSpan creates a new tracing span.
func (o *OTelProvider) StartSpan(ctx context.Context, name string, attributes map[string]string) (context.Context, core.Span) {
	ctx, span := o.tracer.Start(ctx, name)

	if len(attributes) > 0 {
		var attrs []attribute.KeyValue
		for k, v := range attributes {
			attrs = append(attrs, attribute.String(k, v))
		}
		span.SetAttributes(attrs...)
	}

	return ctx, &otelSpan{span: span}
}

// otelSpan implements the core.Span interface.
type otelSpan struct {
	span trace.Span
}

func (s *otelSpan) End() {
	s.span.End()
}

func (s *otelSpan) RecordError(err error) {
	s.span.RecordError(err)
}

func (s *otelSpan) SetAttributes(attributes map[string]string) {
	var attrs []attribute.KeyValue
	for k, v := range attributes {
		attrs = append(attrs, attribute.String(k, v))
	}
	s.span.SetAttributes(attrs...)
}

// PrometheusMetrics implements the core.MetricsRecorder interface.
type PrometheusMetrics struct {
	counters   map[string]prometheus.Counter
	histograms map[string]prometheus.Histogram
	gauges     map[string]prometheus.Gauge
	// In production, sync.Map or mutexes should protect these maps if dynamic metric creation is needed,
	// but static registration via NewPrometheusMetrics is preferred.
}

// NewPrometheusMetrics creates a generic metrics recorder.
// Note: It's usually better to pre-register specific metrics. This naive dynamic map
// is simplified for scaffolding. A strict implementation would inject concrete structs.
func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		counters:   make(map[string]prometheus.Counter),
		histograms: make(map[string]prometheus.Histogram),
		gauges:     make(map[string]prometheus.Gauge),
	}
}

// RegisterCounter registers a counter to avoid race conditions during map access.
func (p *PrometheusMetrics) RegisterCounter(name, help string) {
	p.counters[name] = promauto.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})
}

func (p *PrometheusMetrics) IncCounter(ctx context.Context, name string, labels map[string]string) {
	if c, ok := p.counters[name]; ok {
		c.Inc()
	}
}

func (p *PrometheusMetrics) ObserveHistogram(ctx context.Context, name string, value float64, labels map[string]string) {
	if h, ok := p.histograms[name]; ok {
		h.Observe(value)
	}
}

func (p *PrometheusMetrics) SetGauge(ctx context.Context, name string, value float64, labels map[string]string) {
	if g, ok := p.gauges[name]; ok {
		g.Set(value)
	}
}
