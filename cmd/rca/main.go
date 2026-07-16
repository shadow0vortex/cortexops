package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shadow0vortex/cortexops/internal/rca/llm"
	"github.com/shadow0vortex/cortexops/internal/rca/memory"
	"github.com/shadow0vortex/cortexops/internal/rca/rag"
	"github.com/shadow0vortex/cortexops/pkg/broker"
	"github.com/shadow0vortex/cortexops/pkg/logger"
	"github.com/shadow0vortex/cortexops/pkg/telemetry"
	correlationv1 "github.com/shadow0vortex/cortexops/api/v1"
	"google.golang.org/protobuf/proto"
)

func main() {
	log := logger.New(logger.Config{Level: "info"})
	slog.SetDefault(log)

	log.Info("Starting CortexOps RCA Service")

	// Setup Diagnostics Server
	http.HandleFunc("/debug/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":  "ok",
			"service": "rca",
		})
	})
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Info("Starting diagnostics server on :9091")
		if err := http.ListenAndServe(":9091", nil); err != nil {
			log.Error("Diagnostics server failed", "error", err)
		}
	}()

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

	// Initialize RCA Stream
	err = natsBroker.InitStream("RCA", []string{"cortex.rca.>"})
	if err != nil {
		log.Error("Failed to initialize RCA stream", "error", err)
		os.Exit(1)
	}

	// Initialize Vector Store
	qdrantURL := os.Getenv("QDRANT_URL")
	if qdrantURL == "" {
		qdrantURL = "http://localhost:6333"
	}
	vectorDB := memory.NewQdrantClient(qdrantURL, "incidents")

	// Initialize Real AI Client
	openAIKey := os.Getenv("OPENAI_API_KEY")
	aiClient := llm.NewOpenAIClient(openAIKey)

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
		log.Info("Validation: RCA processing incident", "incidentID", incident.IncidentId)
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
