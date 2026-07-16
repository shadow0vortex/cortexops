package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/shadow0vortex/cortexops/api/v1"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("expected 'duplicate-storm', 'workflow-idempotency', or 'failure-injection' subcommands")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "duplicate-storm":
		size := 100
		if len(os.Args) > 2 {
			fmt.Sscanf(os.Args[2], "%d", &size)
		}
		runDuplicateStorm(size)
	case "workflow-idempotency":
		runWorkflowIdempotency()
	case "failure-injection":
		scenario := "rollout-fail"
		if len(os.Args) > 2 {
			scenario = os.Args[2]
		}
		runFailureInjection(scenario)
	default:
		fmt.Println("expected 'duplicate-storm', 'workflow-idempotency', or 'failure-injection' subcommands")
		os.Exit(1)
	}
}

func runDuplicateStorm(stormSize int) {
	fmt.Println("Running duplicate-storm chaos test...")
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://admin:cortexpassword@localhost:4222"
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	payloadData := map[string]interface{}{
		"severity": "WARNING",
		"message":  "Simulated deployment failure for chaos testing",
		"context": map[string]string{
			"namespace": "cortexops-demo",
			"pod":       "demo-frontend-x",
		},
	}
	payloadBytes, _ := json.Marshal(payloadData)

	// Base event
	event := &apiv1.TelemetryEvent{
		Id:        "chaos-duplicate-12345",
		Source:    "chaos-runner",
		Timestamp: timestamppb.Now(),
		Type:      "k8s.deployment.warning",
		Payload:   payloadBytes,
	}
	data, err := proto.Marshal(event)
	if err != nil {
		log.Fatalf("Failed to marshal: %v", err)
	}

	subject := "cortex.telemetry.k8s.WARNING"

	fmt.Printf("Injecting %d identical telemetry events to %s...\n", stormSize, subject)
	start := time.Now()
	for i := 0; i < stormSize; i++ {
		// Use exact same message to test deduplication
		err := nc.Publish(subject, data)
		if err != nil {
			log.Fatalf("Failed to publish: %v", err)
		}
	}
	nc.Flush()
	fmt.Printf("Publishing and flushing %d events took %v\n", stormSize, time.Since(start))

	fmt.Println("Waiting for correlation window (5s)...")
	time.Sleep(5 * time.Second)

	fmt.Println("Verifying incident count...")
	resp, err := http.Get("http://localhost:9091/debug/incidents/active")
	if err != nil {
		log.Fatalf("Failed to get incidents: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var incidents map[string]interface{}
	json.Unmarshal(body, &incidents)

	// Since we send the exact same ID and the correlator window groups by signature
	fmt.Printf("Active incidents response: %s\n", string(body))
	fmt.Println("Duplicate storm handled. If there is only 1 incident generated from this burst, deduplication is successful.")
}

func runWorkflowIdempotency() {
	fmt.Println("Running workflow-idempotency chaos test...")
	natsURL := os.Getenv("NATS_URL")
	if natsURL == "" {
		natsURL = "nats://admin:cortexpassword@localhost:4222"
	}
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Send duplicate RCA reports which trigger remediation
	report := &apiv1.RCAReport{
		RcaId:       "chaos-rca-999",
		IncidentId:  "chaos-incident-999",
		Analysis:    "Simulated chaos idempotency root cause",
		ConfidenceScore: 0.99,
		AdvisoryRemediationSteps: []string{"restart-pod"},
		GeneratedAt: timestamppb.Now(),
	}
	data, err := proto.Marshal(report)
	if err != nil {
		log.Fatalf("Failed to marshal RCA report: %v", err)
	}

	fmt.Println("Injecting 5 identical RCA reports to cortex.rca.report...")
	for i := 0; i < 5; i++ {
		nc.Publish("cortex.rca.report", data)
	}
	nc.Flush()

	fmt.Println("Waiting for Remediation service to trigger workflows (5s)...")
	time.Sleep(5 * time.Second)

	fmt.Println("Workflow idempotency test completed. Verify in Temporal UI that only ONE workflow was created for incident 'chaos-incident-999'.")
}

func runFailureInjection(scenario string) {
	fmt.Printf("Running failure-injection: %s...\n", scenario)
	cmd := exec.Command("go", "run", "./cmd/demo/main.go", "inject", "-scenario="+scenario)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failure injection failed: %v", err)
	}

	fmt.Println("Waiting 10 seconds for events to propagate and workflows to complete...")
	time.Sleep(10 * time.Second)
	
	resp, err := http.Get("http://localhost:9091/debug/incidents/active")
	if err == nil {
		defer resp.Body.Close()
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("Post-injection incidents: %s\n", string(body))
	}
}
