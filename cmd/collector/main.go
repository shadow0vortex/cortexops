package main

import (
	"context"
	"encoding/json"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
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
	flag.Parse()

	log := logger.New(logger.Config{Level: "info"})
	slog.SetDefault(log)

	log.Info("Starting CortexOps Collector")

	// Setup Diagnostics Server
	k8sConnected := false
	http.HandleFunc("/debug/healthz", func(w http.ResponseWriter, r *http.Request) {
		status := "ok"
		if !k8sConnected {
			status = "degraded"
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"status":        status,
			"k8s_connected": k8sConnected,
			"service":       "collector",
		})
	})
	go func() {
		log.Info("Starting diagnostics server on :9091")
		if err := http.ListenAndServe(":9091", nil); err != nil {
			log.Error("Diagnostics server failed", "error", err)
		}
	}()

	// Initialize K8s client
	var config *rest.Config
	var err error
	var client kubernetes.Interface

	config, err = rest.InClusterConfig()
	if err != nil {
		log.Warn("Failed to get in-cluster config, falling back to local kubeconfig", "error", err)
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Error("Failed to load kubeconfig", "error", err)
		} else if os.Getenv("DOCKER_COMPOSE_ENV") == "true" {
			log.Info("DOCKER_COMPOSE_ENV=true detected, rewriting loopback to host.docker.internal")
			config.Host = strings.ReplaceAll(config.Host, "127.0.0.1", "host.docker.internal")
			config.Host = strings.ReplaceAll(config.Host, "localhost", "host.docker.internal")
			config.Insecure = true
			config.CAFile = ""
			config.CAData = nil
		}
	}

	if config != nil {
		client, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Error("Failed to create kubernetes client", "error", err)
		} else {
			version, err := client.Discovery().ServerVersion()
			if err != nil {
				log.Error("Failed to ping Kubernetes cluster", "error", err)
			} else {
				k8sConnected = true
				log.Info("Connected to Kubernetes cluster",
					"ClusterName", config.Host,
					"API_Server", config.Host,
					"NamespaceScope", "All",
					"Version", version.String())
			}
		}
	}

	if !k8sConnected {
		log.Warn("CRITICAL: Running in degraded mode. Fake clientset fallback is DISABLED. Event ingestion is stopped. Check /debug/healthz")
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

	// Initialize Streams
	err = natsBroker.InitStream("TELEMETRY", []string{"cortex.telemetry.>"})
	if err != nil {
		log.Error("Failed to initialize telemetry stream", "error", err)
		os.Exit(1)
	}

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
	if k8sConnected {
		go func() {
			if err := watcher.Start(ctx); err != nil {
				log.Error("Watcher failed", "error", err)
			}
		}()
	} else {
		log.Warn("Skipping watcher start due to Kubernetes disconnection.")
	}

	log.Info("Collector is running. Press Ctrl+C to stop.")
	<-ctx.Done()
	log.Info("Shutting down collector...")
}
