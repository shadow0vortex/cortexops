package telemetry

import (
	"context"
	"os"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/shadow0vortex/cortexops/pkg/core"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
)

// OTelProvider implements the core.TraceProvider interface.
type OTelProvider struct {
	tracer trace.Tracer
}

// NewOTelProvider creates a new trace provider for a specific service.
// It wires up the OTLP exporter if OTLP_ENDPOINT is set.
func NewOTelProvider(serviceName string) *OTelProvider {
	endpoint := os.Getenv("OTLP_ENDPOINT")
	if endpoint != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		exp, err := otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(endpoint),
			otlptracegrpc.WithInsecure(), // Adjust for production TLS
		)
		if err == nil {
			res, _ := resource.Merge(
				resource.Default(),
				resource.NewWithAttributes(
					semconv.SchemaURL,
					semconv.ServiceName(serviceName),
				),
			)
			tp := sdktrace.NewTracerProvider(
				sdktrace.WithBatcher(exp),
				sdktrace.WithResource(res),
			)
			otel.SetTracerProvider(tp)
			otel.SetTextMapPropagator(propagation.TraceContext{})
		}
	}

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
	mu         sync.RWMutex
	counters   map[string]prometheus.Counter
	histograms map[string]prometheus.Histogram
	gauges     map[string]prometheus.Gauge
}

// NewPrometheusMetrics creates a generic metrics recorder.
func NewPrometheusMetrics() *PrometheusMetrics {
	return &PrometheusMetrics{
		counters:   make(map[string]prometheus.Counter),
		histograms: make(map[string]prometheus.Histogram),
		gauges:     make(map[string]prometheus.Gauge),
	}
}

// RegisterCounter registers a counter to avoid race conditions during map access.
func (p *PrometheusMetrics) RegisterCounter(name, help string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.counters[name] = promauto.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: help,
	})
}

// RegisterGauge registers a gauge.
func (p *PrometheusMetrics) RegisterGauge(name, help string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.gauges[name] = promauto.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: help,
	})
}

func (p *PrometheusMetrics) IncCounter(ctx context.Context, name string, labels map[string]string) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if c, ok := p.counters[name]; ok {
		c.Inc()
	}
}

func (p *PrometheusMetrics) ObserveHistogram(ctx context.Context, name string, value float64, labels map[string]string) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if h, ok := p.histograms[name]; ok {
		h.Observe(value)
	}
}

func (p *PrometheusMetrics) SetGauge(ctx context.Context, name string, value float64, labels map[string]string) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if g, ok := p.gauges[name]; ok {
		g.Set(value)
	}
}
