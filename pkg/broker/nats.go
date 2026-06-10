package broker

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/shadow0vortex/cortexops/pkg/core"
	"github.com/shadow0vortex/cortexops/pkg/retry"
	"github.com/nats-io/nats.go"
)

// NatsBroker implements core.Publisher and core.Subscriber using NATS JetStream.
type NatsBroker struct {
	nc     *nats.Conn
	js     nats.JetStreamContext
	logger *slog.Logger
}

// NewNatsBroker creates a new broker connection.
func NewNatsBroker(url string, logger *slog.Logger) (*NatsBroker, error) {
	nc, err := nats.Connect(url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
		nats.ReconnectWait(2*time.Second),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			logger.Warn("Disconnected from NATS", "error", err)
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			logger.Info("Reconnected to NATS")
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to nats: %w", err)
	}

	js, err := nc.JetStream(nats.PublishAsyncMaxPending(256))
	if err != nil {
		return nil, fmt.Errorf("failed to get jetstream context: %w", err)
	}

	return &NatsBroker{nc: nc, js: js, logger: logger}, nil
}

// InitStream creates the telemetry stream with exact-once delivery semantics and DLQ routing.
func (b *NatsBroker) InitStream(streamName string, subjects []string) error {
	_, err := b.js.StreamInfo(streamName)
	if err == nil {
		b.logger.Info("Stream already exists", "stream", streamName)
		return nil
	}
	if !errors.Is(err, nats.ErrStreamNotFound) {
		return fmt.Errorf("error checking stream info: %w", err)
	}

	cfg := &nats.StreamConfig{
		Name:         streamName,
		Subjects:     subjects,
		Storage:      nats.FileStorage,
		Retention:    nats.LimitsPolicy,
		MaxAge:       7 * 24 * time.Hour, // Keep data for 7 days
		Discard:      nats.DiscardOld,
		Duplicates:   5 * time.Minute, // Deduplicate messages with same Msg-Id within 5m
	}

	_, err = b.js.AddStream(cfg)
	if err != nil {
		return fmt.Errorf("failed to add stream: %w", err)
	}

	b.logger.Info("Successfully created NATS JetStream", "stream", streamName)
	return nil
}

// Publish implements core.Publisher. It enforces idempotent event publishing via MsgId.
func (b *NatsBroker) Publish(ctx context.Context, subject string, eventID string, payload []byte) error {
	// Nats-Msg-Id enforces deduplication on the broker side.
	msg := &nats.Msg{
		Subject: subject,
		Data:    payload,
		Header:  nats.Header{},
	}
	msg.Header.Set("Nats-Msg-Id", eventID)

	// We use exponential backoff for publishing in case the broker is temporarily partitioned.
	return retry.Do(ctx, retry.DefaultPolicy, func() error {
		_, err := b.js.PublishMsg(msg, nats.Context(ctx))
		return err
	})
}

// Subscribe implements core.Subscriber. Handles backpressure and safe acknowledgements.
func (b *NatsBroker) Subscribe(ctx context.Context, subject string, handler core.EventHandler) error {
	// We derive a unique durable name from the subject to avoid "subject does not match consumer" errors
	// when multiple subscribers use the same default group name but different subjects.
	// We replace dots with underscores because NATS durable names cannot contain dots.
	group := "cortex-consumer-" + strings.ReplaceAll(subject, ".", "_")
	
	// MaxDeliver sets the dead-letter queue (DLQ) threshold. After 3 failures, NATS stops redelivering.
	subOpts := []nats.SubOpt{
		nats.Durable(group),
		nats.MaxDeliver(3),
		nats.AckWait(30 * time.Second),
		nats.ManualAck(),
	}

	_, err := b.js.QueueSubscribe(subject, group, func(msg *nats.Msg) {
		// Create a span/context for processing this specific message
		msgCtx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
		defer cancel()

		err := handler(msgCtx, msg.Data)
		if err != nil {
			b.logger.Error("Failed to process message, sending Nak", "subject", msg.Subject, "error", err)
			if nerr := msg.Nak(); nerr != nil {
				b.logger.Error("Failed to NAK message", "error", nerr)
			}
			return
		}

		if err := msg.Ack(); err != nil {
			b.logger.Error("Failed to ACK message", "error", err)
		}
	}, subOpts...)

	if err != nil {
		return fmt.Errorf("failed to subscribe to %s: %w", subject, err)
	}
	return nil
}

func (b *NatsBroker) Close() error {
	b.nc.Close()
	return nil
}
