package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/shadow0vortex/cortexops/internal/diagnostics"
	"github.com/shadow0vortex/cortexops/internal/topology/discovery"
	"github.com/shadow0vortex/cortexops/internal/topology/graph"
	"github.com/shadow0vortex/cortexops/pkg/logger"
	"github.com/shadow0vortex/cortexops/pkg/telemetry"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	log := logger.New(logger.Config{Level: "info"})
	slog.SetDefault(log)

	log.Info("Starting CortexOps Topology Service")

	// Initialize K8s client
	var config *rest.Config
	var err error
	config, err = rest.InClusterConfig()
	if err != nil {
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

	// Initialize metrics
	metrics := telemetry.NewPrometheusMetrics()

	// Initialize Graph Store
	graphStore := graph.NewMemoryGraphStore()

	// Initialize and Start Discovery
	discoverer := discovery.NewK8sDiscovery(client, graphStore, metrics, log)
	
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Info("Starting Topology Discovery Engine")
	go func() {
		if err := discoverer.Start(ctx); err != nil {
			log.Error("Discovery engine failed", "error", err)
		}
	}()

	// Initialize and Start Diagnostics API
	diagAPI := diagnostics.NewAPI(graphStore)
	mux := http.NewServeMux()
	diagAPI.RegisterRoutes(mux)

	port := os.Getenv("DIAG_PORT")
	if port == "" {
		port = "9091"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Info("Starting Diagnostics API", "port", port)
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("Diagnostics API failed", "error", err)
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down topology service...")
	
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("Server shutdown failed", "error", err)
	}
}
