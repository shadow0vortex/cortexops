package chaos

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

// Injector manages programmatic fault injection against Docker Compose services.
// It operates against the compose network created by `make devup`.
type Injector struct {
	// networkName is the Docker Compose network name (e.g., "compose_default").
	networkName string
}

// NewInjector creates a new chaos injector targeting the given Docker network.
// If networkName is empty, it defaults to "compose_default".
func NewInjector(networkName string) *Injector {
	if networkName == "" {
		networkName = "compose_default"
	}
	return &Injector{networkName: networkName}
}

// PartitionNetwork disconnects a container from the Docker network for the specified duration,
// then automatically reconnects it. This simulates a network partition to a dependency
// such as "compose-postgres-1" or "compose-qdrant-1".
func (i *Injector) PartitionNetwork(ctx context.Context, containerName string, duration time.Duration) error {
	fmt.Printf("[CHAOS] Partitioning %s from network %s for %v\n", containerName, i.networkName, duration)

	// Disconnect the container from the compose network.
	cmd := exec.CommandContext(ctx, "docker", "network", "disconnect", i.networkName, containerName)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to disconnect %s: %w (output: %s)", containerName, err, strings.TrimSpace(string(out)))
	}

	// Schedule automatic reconnection after the specified duration.
	go func() {
		select {
		case <-time.After(duration):
		case <-ctx.Done():
		}

		reconnectCmd := exec.Command("docker", "network", "connect", i.networkName, containerName)
		if out, err := reconnectCmd.CombinedOutput(); err != nil {
			fmt.Printf("[CHAOS] WARNING: failed to reconnect %s: %v (output: %s)\n", containerName, err, strings.TrimSpace(string(out)))
		} else {
			fmt.Printf("[CHAOS] Restored %s to network %s\n", containerName, i.networkName)
		}
	}()

	return nil
}

// RestoreNetwork manually reconnects a container to the network.
// Use this in test cleanup to guarantee restoration even if the timed goroutine hasn't fired.
func (i *Injector) RestoreNetwork(containerName string) error {
	cmd := exec.Command("docker", "network", "connect", i.networkName, containerName)
	if out, err := cmd.CombinedOutput(); err != nil {
		// Ignore "already connected" errors during cleanup.
		if strings.Contains(string(out), "already exists") {
			return nil
		}
		return fmt.Errorf("failed to reconnect %s: %w (output: %s)", containerName, err, strings.TrimSpace(string(out)))
	}
	return nil
}

// KillContainer sends SIGKILL to a running container, simulating an OOMKill or crash.
func (i *Injector) KillContainer(ctx context.Context, containerName string) error {
	fmt.Printf("[CHAOS] Killing container %s\n", containerName)

	cmd := exec.CommandContext(ctx, "docker", "kill", containerName)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to kill %s: %w (output: %s)", containerName, err, strings.TrimSpace(string(out)))
	}
	return nil
}

// RestartContainer restarts a previously killed container.
func (i *Injector) RestartContainer(ctx context.Context, containerName string) error {
	fmt.Printf("[CHAOS] Restarting container %s\n", containerName)

	cmd := exec.CommandContext(ctx, "docker", "start", containerName)
	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to restart %s: %w (output: %s)", containerName, err, strings.TrimSpace(string(out)))
	}
	return nil
}

// InjectAPILatency adds network latency to a container using `tc` via docker exec.
// This requires the container to have the `tc` (iproute2) tool available.
func (i *Injector) InjectAPILatency(ctx context.Context, containerName string, latency time.Duration, duration time.Duration) error {
	latencyMs := latency.Milliseconds()
	fmt.Printf("[CHAOS] Injecting %dms latency into %s for %v\n", latencyMs, containerName, duration)

	// Add latency using tc inside the container.
	addCmd := exec.CommandContext(ctx, "docker", "exec", containerName,
		"tc", "qdisc", "add", "dev", "eth0", "root", "netem", "delay", fmt.Sprintf("%dms", latencyMs))
	if out, err := addCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to inject latency into %s: %w (output: %s)", containerName, err, strings.TrimSpace(string(out)))
	}

	// Schedule automatic removal.
	go func() {
		select {
		case <-time.After(duration):
		case <-ctx.Done():
		}

		removeCmd := exec.Command("docker", "exec", containerName,
			"tc", "qdisc", "del", "dev", "eth0", "root")
		if out, err := removeCmd.CombinedOutput(); err != nil {
			fmt.Printf("[CHAOS] WARNING: failed to remove latency from %s: %v (output: %s)\n", containerName, err, strings.TrimSpace(string(out)))
		} else {
			fmt.Printf("[CHAOS] Removed latency injection from %s\n", containerName)
		}
	}()

	return nil
}
