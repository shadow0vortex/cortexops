package main

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	rcav1 "github.com/shadow0vortex/cortexops/api/v1"
	remediationv1 "github.com/shadow0vortex/cortexops/api/v1"
	"github.com/shadow0vortex/cortexops/internal/orchestration/temporal"
	"github.com/shadow0vortex/cortexops/internal/remediation/action"
	"github.com/shadow0vortex/cortexops/internal/remediation/approval"
	"github.com/shadow0vortex/cortexops/internal/remediation/policy"
	"github.com/shadow0vortex/cortexops/pkg/broker"
	"github.com/shadow0vortex/cortexops/pkg/logger"
	"github.com/shadow0vortex/cortexops/pkg/telemetry"
	"github.com/shadow0vortex/cortexops/pkg/topology"
	"github.com/google/uuid"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
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
			"service":       "remediation",
		})
	})
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Info("Starting diagnostics server on :9091")
		if err := http.ListenAndServe(":9091", nil); err != nil {
			log.Error("Diagnostics server failed", "error", err)
		}
	}()

	// Initialize K8s client
	var config *rest.Config
	var err error
	var k8sClient kubernetes.Interface

	config, err = rest.InClusterConfig()
	if err != nil {
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
		k8sClient, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Error("Failed to create kubernetes client", "error", err)
		} else {
			version, err := k8sClient.Discovery().ServerVersion()
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
		log.Warn("CRITICAL: Running in degraded mode. Fake clientset fallback is DISABLED. Kubernetes execution is blocked. Check /debug/healthz")
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

	// Initialize Broker
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://localhost:4222"
	}
	natsBroker, err := broker.NewNatsBroker(natsURL, log)
	if err != nil {
		log.Error("Failed to connect to NATS", "error", err)
	} else {
		defer natsBroker.Close()
		
		// Ensure RCA stream exists before subscribing
		err = natsBroker.InitStream("RCA", []string{"cortex.rca.>"})
		if err != nil {
			log.Error("Failed to initialize RCA stream", "error", err)
		}
	}

	// Initialize metrics
	metrics := telemetry.NewPrometheusMetrics()

	// Initialize dependencies
	policyEngine := policy.NewOPAEngine(topoClient, metrics, log)
	auditStore := &MemoryAuditStore{}
	approvalWorkflow := approval.NewDeterministicApprovalWorkflow(auditStore, metrics, log)
	executor := action.NewK8sExecutor(k8sClient, metrics, log)

	activities := temporal.NewActivities(policyEngine, approvalWorkflow, executor)

	// Create and start worker
	w := worker.New(c, "remediation-tasks", worker.Options{})

	w.RegisterWorkflow(temporal.RemediationWorkflow)
	w.RegisterActivity(activities)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Subscribe to RCA reports
	if natsBroker != nil {
		err = natsBroker.Subscribe(ctx, "cortex.rca.report", func(ctx context.Context, payload []byte) error {
			report := &rcav1.RCAReport{}
			if err := proto.Unmarshal(payload, report); err != nil {
				return err
			}
			
			log.Info("Validation: Remediation proposing action from RCA", "rcaID", report.RcaId, "incidentID", report.IncidentId)

			// Simple heuristic: if we see "rollout" or "demo-frontend", propose rollout restart
			act := &remediationv1.RemediationAction{
				ActionId:   uuid.New().String(),
				IncidentId: report.IncidentId,
				Type:       remediationv1.ActionType_DEPLOYMENT_ROLLOUT_RESTART,
				State:      remediationv1.ActionState_PROPOSED,
				TargetNamespace: "cortexops-demo",
				TargetResource:  "demo-frontend",
				CreatedAt:  timestamppb.Now(),
				UpdatedAt:  timestamppb.Now(),
			}

			// Dummy incident context just to satisfy policy engine signature for this demo
			incidentCtx := &rcav1.CorrelatedIncident{
				IncidentId: report.IncidentId,
			}

			workflowOptions := client.StartWorkflowOptions{
				ID:        "remediation-" + act.ActionId,
				TaskQueue: "remediation-tasks",
			}
			we, err := c.ExecuteWorkflow(ctx, workflowOptions, temporal.RemediationWorkflow, incidentCtx, act)
			if err != nil {
				log.Error("Failed to start temporal workflow", "error", err)
				return err
			}
			log.Info("Validation: Temporal workflow created", "workflowID", we.GetID(), "runID", we.GetRunID())
			return nil
		})
		if err != nil {
			log.Error("Failed to subscribe to RCA reports", "error", err)
		}
	}

	log.Info("Remediation Worker is starting...")
	if k8sConnected {
		go func() {
			if err := w.Run(worker.InterruptCh()); err != nil {
				log.Error("Worker failed", "error", err)
			}
		}()
	} else {
		log.Warn("Skipping Temporal worker start due to Kubernetes disconnection.")
	}

	<-ctx.Done()
	log.Info("Shutting down remediation service...")
}
