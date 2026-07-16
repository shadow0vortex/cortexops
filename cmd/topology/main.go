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

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/shadow0vortex/cortexops/internal/diagnostics"
	"github.com/shadow0vortex/cortexops/internal/topology/discovery"
	"github.com/shadow0vortex/cortexops/internal/topology/graph"
	"github.com/shadow0vortex/cortexops/pkg/logger"
	"github.com/shadow0vortex/cortexops/pkg/telemetry"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
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
	var client kubernetes.Interface

	config, err = rest.InClusterConfig()
	if err != nil {
		kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
		if err != nil {
			log.Warn("Failed to load kubeconfig, falling back to fake clientset for development", "error", err)
			client = fake.NewSimpleClientset()
		}
	}

	if client == nil {
		client, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Error("Failed to create kubernetes client", "error", err)
			os.Exit(1)
		}
	}

	// Initialize metrics
	metrics := telemetry.NewPrometheusMetrics()

	// Initialize Graph Store
	graphStore := graph.NewMemoryGraphStore()

	// Initialize Persister
	pgURL := os.Getenv("POSTGRES_URL")
	if pgURL != "" {
		persister, err := graph.NewGraphPersister(pgURL, graphStore, log)
		if err != nil {
			log.Warn("Failed to initialize graph persister", "error", err)
		} else {
			defer persister.Close()
			if err := persister.Restore(context.Background()); err != nil {
				log.Error("Failed to restore topology graph from snapshot", "error", err)
			}
			
			go func() {
				ticker := time.NewTicker(15 * time.Second)
				defer ticker.Stop()
				for {
					select {
					case <-context.Background().Done():
						return
					case <-ticker.C:
						if err := persister.SaveAsync(context.Background()); err != nil {
							log.Error("Failed to save async snapshot", "error", err)
						}
					}
				}
			}()
		}
	}

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

	diagAPI := diagnostics.NewAPI(graphStore)
	mux := http.NewServeMux()
	diagAPI.RegisterRoutes(mux)
	mux.Handle("/metrics", promhttp.Handler())

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
