package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shadow0vortex/cortexops/internal/correlation/causal"
	"github.com/shadow0vortex/cortexops/internal/correlation/engine"
	"github.com/shadow0vortex/cortexops/internal/correlation/heuristics"
	"github.com/shadow0vortex/cortexops/pkg/broker"
	"github.com/shadow0vortex/cortexops/pkg/logger"
	"github.com/shadow0vortex/cortexops/pkg/telemetry"
	"github.com/shadow0vortex/cortexops/pkg/topology"
	eventsv1 "github.com/shadow0vortex/cortexops/api/v1"
	"google.golang.org/protobuf/proto"
)

func main() {
	log := logger.New(logger.Config{Level: "info"})
	slog.SetDefault(log)

	log.Info("Starting CortexOps Correlator")

	// Setup Diagnostics Server
	http.HandleFunc("/debug/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":  "ok",
			"service": "correlator",
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

	// Initialize Streams
	err = natsBroker.InitStream("INCIDENTS", []string{"cortex.incident.>"})
	if err != nil {
		log.Error("Failed to initialize incidents stream", "error", err)
		os.Exit(1)
	}

	// Initialize Topology Client
	topoURL := os.Getenv("TOPOLOGY_URL")
	if topoURL == "" {
		topoURL = "http://topology:9091"
	}
	topoClient := topology.NewHTTPClient(topoURL)

	// Initialize metrics
	metrics := telemetry.NewPrometheusMetrics()

	// Initialize Engine
	scorer := heuristics.NewScorer(topoClient)
	chainBuilder := causal.NewChainBuilder(topoClient)
	corrEngine := engine.NewEngine(scorer, chainBuilder, natsBroker, metrics, log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Subscribe to telemetry
	subjects := []string{"cortex.telemetry.k8s.INFO", "cortex.telemetry.k8s.WARNING"}
	for _, sub := range subjects {
		err := natsBroker.Subscribe(ctx, sub, func(ctx context.Context, payload []byte) error {
			envelope := &eventsv1.TelemetryEnvelope{}
			if err := proto.Unmarshal(payload, envelope); err != nil {
				return err
			}
			log.Info("Validation: Correlator processing telemetry", "eventID", envelope.EventId, "subject", sub)
			return corrEngine.ProcessEvent(ctx, envelope)
		})
		if err != nil {
			log.Error("Failed to subscribe to telemetry", "subject", sub, "error", err)
			os.Exit(1)
		}
	}

	// Periodically flush incidents
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				corrEngine.FlushWindows(ctx)
			}
		}
	}()

	log.Info("Correlator is running.")
	<-ctx.Done()
	log.Info("Shutting down correlator...")
}
