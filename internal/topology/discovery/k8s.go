package discovery

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	topologyv1 "github.com/shadow0vortex/cortexops/api/v1"
	"github.com/shadow0vortex/cortexops/pkg/core"
	"google.golang.org/protobuf/types/known/timestamppb"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// GraphStore interface matching the in-memory graph methods.
// Using localized interface to enforce DI and boundary separation.
type GraphStore interface {
	UpsertNode(ctx context.Context, node *topologyv1.TopologyNode) error
	DeleteNode(ctx context.Context, nodeID string) error
	UpsertEdge(ctx context.Context, sourceID, targetID string, relationship topologyv1.EdgeType) error
}

// K8sDiscovery watches Kubernetes state and deterministically syncs it to the GraphStore.
type K8sDiscovery struct {
	client  kubernetes.Interface
	store   GraphStore
	metrics core.MetricsRecorder
	logger  *slog.Logger
}

func NewK8sDiscovery(client kubernetes.Interface, store GraphStore, metrics core.MetricsRecorder, logger *slog.Logger) *K8sDiscovery {
	return &K8sDiscovery{
		client:  client,
		store:   store,
		metrics: metrics,
		logger:  logger,
	}
}

// Start begins the informer reconciliation loops.
func (d *K8sDiscovery) Start(ctx context.Context) error {
	factory := informers.NewSharedInformerFactory(d.client, 10*time.Minute)
	
	podInformer := factory.Core().V1().Pods().Informer()
	podInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) { d.handlePod(ctx, obj, false) },
		UpdateFunc: func(old, new interface{}) { d.handlePod(ctx, new, false) },
		DeleteFunc: func(obj interface{}) { d.handlePod(ctx, obj, true) },
	})

	// Setup other informers (Deployments, Services, etc.) similarly...

	d.logger.Info("Starting Topology Discovery Engine")
	factory.Start(ctx.Done())
	
	if !cache.WaitForCacheSync(ctx.Done(), podInformer.HasSynced) {
		return fmt.Errorf("failed to sync pod cache")
	}

	<-ctx.Done()
	return nil
}

func (d *K8sDiscovery) handlePod(ctx context.Context, obj interface{}, isDelete bool) {
	pod, ok := obj.(*corev1.Pod)
	if !ok {
		return
	}

	nodeID := fmt.Sprintf("pod/%s/%s", pod.Namespace, pod.Name)

	if isDelete {
		d.store.DeleteNode(ctx, nodeID)
		d.metrics.IncCounter(ctx, "cortexops_topology_nodes_total", map[string]string{"type": "POD", "action": "deleted"})
		return
	}

	node := &topologyv1.TopologyNode{
		Id:          nodeID,
		Type:        topologyv1.NodeType_POD,
		Namespace:   pod.Namespace,
		Name:        pod.Name,
		Labels:      pod.Labels,
		LastUpdated: timestamppb.Now(),
	}

	d.store.UpsertNode(ctx, node)
	d.metrics.SetGauge(ctx, "cortexops_topology_nodes_total", 1, map[string]string{"type": "POD"})

	// Create edge: POD is SCHEDULED_ON NODE
	if pod.Spec.NodeName != "" {
		targetNodeID := fmt.Sprintf("node/%s", pod.Spec.NodeName)
		d.store.UpsertEdge(ctx, nodeID, targetNodeID, topologyv1.EdgeType_SCHEDULED_ON)
	}
}
