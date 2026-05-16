package chaos

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

// Injector manages programmatic fault injection against external dependencies.
type Injector struct {
	// Pointers to local dependencies (e.g., Docker container names for local Kube/NATS)
}

func NewInjector() *Injector {
	return &Injector{}
}

// PartitionNetwork simulates a network drop to a target port using iptables/tc or docker network disconnect.
func (i *Injector) PartitionNetwork(ctx context.Context, target string, duration time.Duration) error {
	// Example stub: in a real framework, this would use `tc` or `docker network disconnect`
	// exec.CommandContext(ctx, "docker", "network", "disconnect", "kind", target).Run()
	
	fmt.Printf("[CHAOS] Partitioning network to %s for %v\n", target, duration)
	
	go func() {
		time.Sleep(duration)
		// exec.Command("docker", "network", "connect", "kind", target).Run()
		fmt.Printf("[CHAOS] Restored network to %s\n", target)
	}()
	return nil
}

// InjectAPILatency slows down responses from the target using toxiproxy or tc.
func (i *Injector) InjectAPILatency(ctx context.Context, target string, latency time.Duration, duration time.Duration) error {
	fmt.Printf("[CHAOS] Injecting %v latency into %s for %v\n", latency, target, duration)
	// Implement latency injection (e.g. toxiproxy-cli target latency -d latency)
	return nil
}

// CrashProcess simulates OOMKill or panic by sending SIGKILL to the CortexOps pod.
func (i *Injector) CrashProcess(ctx context.Context, podName string, namespace string) error {
	fmt.Printf("[CHAOS] Crashing pod %s in namespace %s\n", podName, namespace)
	// kubectl delete pod --grace-period=0 --force
	cmd := exec.CommandContext(ctx, "kubectl", "delete", "pod", podName, "-n", namespace, "--grace-period=0", "--force")
	return cmd.Run()
}
