package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/shadow0vortex/cortexops/internal/collector/k8s"
	"github.com/shadow0vortex/cortexops/pkg/broker"
	"github.com/shadow0vortex/cortexops/pkg/logger"
	"github.com/shadow0vortex/cortexops/pkg/telemetry"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	diagnostics := flag.Bool("diagnostics", false, "Run diagnostics and exit")
	flag.Parse()

	log := logger.New(logger.Config{Level: "info"})
	slog.SetDefault(log)

	if *diagnostics {
		log.Info("Running diagnostics...")
		// Placeholder for actual diagnostic logic
		os.Exit(0)
	}

	log.Info("Starting CortexOps Collector")

	// Initialize K8s client
	var config *rest.Config
	var err error
	config, err = rest.InClusterConfig()
	if err != nil {
		log.Warn("Failed to get in-cluster config, falling back to local kubeconfig", "error", err)
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Error("Failed to load kubeconfig", "error", err)
			os.Exit(1)
		}
	}

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error("Failed to create kubernetes client", "error", err)
		os.Exit(1)
	}

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

	// Initialize metrics
	metrics := telemetry.NewPrometheusMetrics()
	metrics.RegisterCounter("cortexops_telemetry_ingested_total", "Total telemetry ingested")
	metrics.RegisterCounter("cortexops_telemetry_dropped_total", "Total telemetry dropped")
	metrics.RegisterCounter("cortexops_telemetry_published_total", "Total telemetry published")

	// Create and start watcher
	watcher := k8s.NewWatcher(client, natsBroker, metrics, log)

	// Wait for termination
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Info("Collector is starting event watcher...")
	go func() {
		if err := watcher.Start(ctx); err != nil {
			log.Error("Watcher failed", "error", err)
		}
	}()

	log.Info("Collector is running. Press Ctrl+C to stop.")
	<-ctx.Done()
	log.Info("Shutting down collector...")
}
