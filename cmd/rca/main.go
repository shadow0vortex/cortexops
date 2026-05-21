package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/shadow0vortex/cortexops/internal/rca/rag"
	"github.com/shadow0vortex/cortexops/internal/rca/memory"
	"github.com/shadow0vortex/cortexops/pkg/broker"
	"github.com/shadow0vortex/cortexops/pkg/logger"
	"github.com/shadow0vortex/cortexops/pkg/telemetry"
	correlationv1 "github.com/shadow0vortex/cortexops/api/v1"
	"google.golang.org/protobuf/proto"
)

// MockAI implements core.LLMClient and core.EmbeddingClient for demo purposes.
type MockAI struct{}

func (m *MockAI) GenerateEmbeddings(ctx context.Context, text string) ([]float32, error) {
	return make([]float32, 384), nil // Fake embedding vector
}

func (m *MockAI) GenerateRCA(ctx context.Context, prompt string) (string, error) {
	return "DEMO ANALYSIS: The incident appears to be caused by a misconfigured image tag in the frontend deployment, leading to ImagePullBackOff.", nil
}

func main() {
	log := logger.New(logger.Config{Level: "info"})
	slog.SetDefault(log)

	log.Info("Starting CortexOps RCA Service")

	// Initialize Broker
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	natsBroker, err := broker.NewNatsBroker(natsURL, log)
	if err != nil {
		log.Error("Failed to connect to NATS", "error", err)
		os.Exit(1)
	}
	defer natsBroker.Close()

	// Initialize Vector Store
	qdrantURL := os.Getenv("QDRANT_URL")
	if qdrantURL == "" {
		qdrantURL = "http://localhost:6333"
	}
	vectorDB := memory.NewQdrantClient(qdrantURL, "incidents")

	// Initialize AI Client (Stubbed for demo)
	aiClient := &MockAI{}

	// Initialize metrics
	metrics := telemetry.NewPrometheusMetrics()

	// Initialize Engine
	rcaEngine := rag.NewEngine(aiClient, aiClient, vectorDB, natsBroker, metrics, log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Subscribe to correlated incidents
	err = natsBroker.Subscribe(ctx, "cortex.incident.correlated", func(ctx context.Context, payload []byte) error {
		incident := &correlationv1.CorrelatedIncident{}
		if err := proto.Unmarshal(payload, incident); err != nil {
			return err
		}
		return rcaEngine.GenerateRCA(ctx, incident)
	})
	if err != nil {
		log.Error("Failed to subscribe to incidents", "error", err)
		os.Exit(1)
	}

	log.Info("RCA Service is running.")
	<-ctx.Done()
	log.Info("Shutting down RCA service...")
}
