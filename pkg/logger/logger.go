package logger

import (
	"context"
	"log/slog"
	"os"
)

// TraceContextKey is the context key for tracing information.
type TraceContextKey string

const (
	TraceIDKey TraceContextKey = "trace_id"
	SpanIDKey  TraceContextKey = "span_id"
)

// Config configures the logger.
type Config struct {
	Level  string
	Format string // "json" or "text"
}

// New creates and returns a configured slog Logger to be injected into services.
// It explicitly avoids mutating the global slog.Default() state.
func New(cfg Config) *slog.Logger {
	var level slog.Level
	switch cfg.Level {
	case "DEBUG", "debug":
		level = slog.LevelDebug
	case "WARN", "warn":
		level = slog.LevelWarn
	case "ERROR", "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level: level,
	}

	var handler slog.Handler
	if cfg.Format == "json" {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	return slog.New(handler)
}

// WithContext extracts trace_id and span_id from context and returns a derived logger.
// It requires the base logger to be explicitly passed, reinforcing dependency injection.
func WithContext(ctx context.Context, baseLogger *slog.Logger) *slog.Logger {
	logger := baseLogger

	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		logger = logger.With("trace_id", traceID)
	}
	if spanID, ok := ctx.Value(SpanIDKey).(string); ok {
		logger = logger.With("span_id", spanID)
	}

	return logger
}
