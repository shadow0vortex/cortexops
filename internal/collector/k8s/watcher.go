package k8s

import (
	"context"
	"log/slog"
	"time"

	"github.com/shadow0vortex/cortexops/internal/collector/normalizer"
	"github.com/shadow0vortex/cortexops/pkg/core"
	"google.golang.org/protobuf/proto"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

// Watcher tails Kubernetes Events using an Informer and publishes them to NATS.
type Watcher struct {
	client     kubernetes.Interface
	publisher  core.Publisher
	metrics    core.MetricsRecorder
	normalizer *normalizer.Normalizer
	logger     *slog.Logger
}

func NewWatcher(client kubernetes.Interface, pub core.Publisher, metrics core.MetricsRecorder, logger *slog.Logger) *Watcher {
	return &Watcher{
		client:     client,
		publisher:  pub,
		metrics:    metrics,
		normalizer: normalizer.New(),
		logger:     logger,
	}
}

// Start begins watching the cluster. It blocks until ctx is canceled.
func (w *Watcher) Start(ctx context.Context) error {
	// Resync period of 10 minutes heals split-brain or missed watch events
	factory := informers.NewSharedInformerFactory(w.client, 10*time.Minute)
	eventInformer := factory.Core().V1().Events().Informer()

	_, _ = eventInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			w.handleEvent(ctx, obj)
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			// K8s sometimes updates the Event count instead of creating new events
			w.handleEvent(ctx, newObj)
		},
	})

	w.logger.Info("Starting Kubernetes Event Watcher")
	factory.Start(ctx.Done())
	
	if !cache.WaitForCacheSync(ctx.Done(), eventInformer.HasSynced) {
		w.logger.Error("Failed to sync kubernetes cache")
		return context.Canceled
	}

	w.logger.Info("Kubernetes informer cache synced")
	<-ctx.Done()
	return nil
}

func (w *Watcher) handleEvent(ctx context.Context, obj interface{}) {
	k8sEvent, ok := obj.(*corev1.Event)
	if !ok {
		w.logger.Error("Received unknown object type in event informer")
		return
	}

	w.metrics.IncCounter(ctx, "cortexops_telemetry_ingested_total", map[string]string{"source": "k8s_event"})

	envelope, err := w.normalizer.NormalizeK8sEvent(ctx, k8sEvent)
	if err != nil {
		w.logger.Error("Failed to normalize K8s event", "error", err)
		w.metrics.IncCounter(ctx, "cortexops_telemetry_dropped_total", map[string]string{"reason": "normalization_failure"})
		return
	}

	payload, err := proto.Marshal(envelope)
	if err != nil {
		w.logger.Error("Failed to marshal protobuf envelope", "error", err)
		w.metrics.IncCounter(ctx, "cortexops_telemetry_dropped_total", map[string]string{"reason": "marshal_failure"})
		return
	}

	subject := "cortex.telemetry.k8s." + envelope.Severity
	err = w.publisher.Publish(ctx, subject, envelope.EventId, payload)
	if err != nil {
		w.logger.Error("Failed to publish telemetry event", "error", err, "event_id", envelope.EventId)
		w.metrics.IncCounter(ctx, "cortexops_telemetry_dropped_total", map[string]string{"reason": "publish_failure"})
		return
	}

	w.metrics.IncCounter(ctx, "cortexops_telemetry_published_total", map[string]string{"subject": subject})
	w.logger.Info("Validation: Collector published telemetry", "eventID", envelope.EventId, "subject", subject)
}
