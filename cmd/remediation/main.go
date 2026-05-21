package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	remediationv1 "github.com/shadow0vortex/cortexops/api/v1"
	"github.com/shadow0vortex/cortexops/internal/orchestration/temporal"
	"github.com/shadow0vortex/cortexops/internal/remediation/action"
	"github.com/shadow0vortex/cortexops/internal/remediation/approval"
	"github.com/shadow0vortex/cortexops/internal/remediation/policy"
	"github.com/shadow0vortex/cortexops/pkg/logger"
	"github.com/shadow0vortex/cortexops/pkg/telemetry"
	"github.com/shadow0vortex/cortexops/pkg/topology"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// MemoryAuditStore is a simple stub for demo purposes.
type MemoryAuditStore struct {
	records []*remediationv1.AuditRecord
}

func (s *MemoryAuditStore) Log(ctx context.Context, record *remediationv1.AuditRecord) error {
	s.records = append(s.records, record)
	return nil
}

func main() {
	log := logger.New(logger.Config{Level: "info"})
	slog.SetDefault(log)

	log.Info("Starting CortexOps Remediation Service")

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
	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error("Failed to create kubernetes client", "error", err)
		os.Exit(1)
	}

	// Initialize Temporal client
	temporalURL := os.Getenv("TEMPORAL_URL")
	if temporalURL == "" {
		temporalURL = "localhost:7233"
	}
	c, err := client.Dial(client.Options{
		HostPort: temporalURL,
	})
	if err != nil {
		log.Error("Failed to connect to Temporal", "error", err)
		os.Exit(1)
	}
	defer c.Close()

	// Initialize Topology Client
	topoURL := os.Getenv("TOPOLOGY_URL")
	if topoURL == "" {
		topoURL = "http://topology:9091"
	}
	topoClient := topology.NewHTTPClient(topoURL)

	// Initialize metrics
	metrics := telemetry.NewPrometheusMetrics()

	// Initialize dependencies
	policyEngine := policy.NewOPAEngine(topoClient, metrics, log)
	auditStore := &MemoryAuditStore{}
	approvalWorkflow := approval.NewSlackWorkflow(auditStore, metrics, log)
	executor := action.NewK8sExecutor(k8sClient, metrics, log)

	activities := temporal.NewActivities(policyEngine, approvalWorkflow, executor)

	// Create and start worker
	w := worker.New(c, "remediation-tasks", worker.Options{})

	w.RegisterWorkflow(temporal.RemediationWorkflow)
	w.RegisterActivity(activities)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	log.Info("Remediation Worker is starting...")
	go func() {
		if err := w.Run(worker.InterruptCh()); err != nil {
			log.Error("Worker failed", "error", err)
		}
	}()

	<-ctx.Done()
	log.Info("Shutting down remediation service...")
}
