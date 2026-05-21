package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	bootstrapCmd := flag.NewFlagSet("bootstrap", flag.ExitOnError)
	injectCmd := flag.NewFlagSet("inject", flag.ExitOnError)
	recoverCmd := flag.NewFlagSet("recover", flag.ExitOnError)

	scenario := injectCmd.String("scenario", "rollout-fail", "Scenario: rollout-fail, crashloop, scaling")

	if len(os.Args) < 2 {
		fmt.Println("expected 'bootstrap', 'inject' or 'recover' subcommands")
		os.Exit(1)
	}

	var config *rest.Config
	var err error
	var clientset kubernetes.Interface

	// Try local kubeconfig
	config, err = clientcmd.BuildConfigFromFlags("", clientcmd.RecommendedHomeFile)
	if err != nil {
		// Fallback to In-Cluster
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("Warning: Failed to load K8s config: %v. Using fake clientset for demo mode.\n", err)
			clientset = fake.NewSimpleClientset()
		}
	}

	if clientset == nil {
		clientset, err = kubernetes.NewForConfig(config)
		if err != nil {
			log.Fatalf("Error creating clientset: %s", err.Error())
		}
	}

	switch os.Args[1] {
	case "bootstrap":
		if err := bootstrapCmd.Parse(os.Args[2:]); err != nil {
			log.Fatalf("failed to parse bootstrap flags: %v", err)
		}
		runBootstrap()
	case "inject":
		if err := injectCmd.Parse(os.Args[2:]); err != nil {
			log.Fatalf("failed to parse inject flags: %v", err)
		}
		runInject(clientset, *scenario)
	case "recover":
		if err := recoverCmd.Parse(os.Args[2:]); err != nil {
			log.Fatalf("failed to parse recover flags: %v", err)
		}
		runRecover(clientset)
	default:
		fmt.Println("expected 'bootstrap', 'inject' or 'recover' subcommands")
		os.Exit(1)
	}
}

func runBootstrap() {
	fmt.Println("Bootstrapping CortexOps Demo Topology...")
	// kubectl also needs to handle no-cluster cases gracefully if called via exec
	cmd := exec.Command("kubectl", "apply", "-f", "sandbox/workloads/demo-topology.yaml")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Printf("Warning: failed to apply topology via kubectl: %v (Is your cluster reachable?)\n", err)
	} else {
		fmt.Println("Topology applied successfully.")
	}
}

func runInject(client kubernetes.Interface, scenario string) {
	fmt.Printf("Injecting failure scenario: %s\n", scenario)
	ctx := context.Background()
	ns := "cortexops-demo"

	switch scenario {
	case "rollout-fail":
		// Set a non-existent image to trigger rollout failure
		patch := []byte(`{"spec": {"template": {"spec": {"containers": [{"name": "web", "image": "nginx:nonexistent-tag"}]}}}}`)
		_, err := client.AppsV1().Deployments(ns).Patch(ctx, "demo-frontend", types.StrategicMergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			log.Fatalf("failed to inject rollout-fail: %v", err)
		}
	case "crashloop":
		// Misconfigure entrypoint to cause crash
		patch := []byte(`{"spec": {"template": {"spec": {"containers": [{"name": "api", "command": ["/bin/false"]}]}}}}`)
		_, err := client.AppsV1().Deployments(ns).Patch(ctx, "demo-api", types.StrategicMergePatchType, patch, metav1.PatchOptions{})
		if err != nil {
			log.Fatalf("failed to inject crashloop: %v", err)
		}
	case "scaling":
		// Trigger scaling pressure
		replicas := int32(10)
		deployment, err := client.AppsV1().Deployments(ns).Get(ctx, "demo-frontend", metav1.GetOptions{})
		if err != nil {
			log.Fatalf("failed to get deployment: %v", err)
		}
		deployment.Spec.Replicas = &replicas
		_, err = client.AppsV1().Deployments(ns).Update(ctx, deployment, metav1.UpdateOptions{})
		if err != nil {
			log.Fatalf("failed to inject scaling pressure: %v", err)
		}
	default:
		log.Fatalf("unknown scenario: %s", scenario)
	}
	fmt.Println("Failure injected.")
}

func runRecover(client kubernetes.Interface) {
	fmt.Println("Recovering from failure scenarios...")
	// Simplest recovery is to re-apply the base topology
	runBootstrap()
}
