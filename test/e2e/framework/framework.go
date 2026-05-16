package framework

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Harness provides an isolated testing environment for E2E validation.
type Harness struct {
	T         *testing.T
	K8s       kubernetes.Interface
	Namespace string
	Ctx       context.Context
	Cancel    context.CancelFunc
}

// Setup creates an isolated namespace and provisions test dependencies.
func Setup(t *testing.T) *Harness {
	t.Helper()
	
	// Defaulting to local kubeconfig for Minikube/Kind CI usage
	config, err := clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		t.Fatalf("Failed to build kubeconfig: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		t.Fatalf("Failed to create k8s client: %v", err)
	}

	ns := fmt.Sprintf("cortexops-test-%s", uuid.New().String()[:8])
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)

	h := &Harness{
		T:         t,
		K8s:       clientset,
		Namespace: ns,
		Ctx:       ctx,
		Cancel:    cancel,
	}

	h.createNamespace()
	t.Cleanup(h.Teardown)
	
	return h
}

func (h *Harness) createNamespace() {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{Name: h.Namespace},
	}
	_, err := h.K8s.CoreV1().Namespaces().Create(h.Ctx, ns, metav1.CreateOptions{})
	if err != nil {
		h.T.Fatalf("Failed to create isolated namespace: %v", err)
	}
}

// Teardown cleans up the isolated environment deterministically.
func (h *Harness) Teardown() {
	h.Cancel() // stop internal contexts
	
	// Use background context for cleanup since test context is canceled
	err := h.K8s.CoreV1().Namespaces().Delete(context.Background(), h.Namespace, metav1.DeleteOptions{})
	if err != nil {
		h.T.Logf("Warning: failed to teardown namespace %s: %v", h.Namespace, err)
	}
}
