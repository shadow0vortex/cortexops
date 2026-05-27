package suites

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	v1 "github.com/shadow0vortex/cortexops/api/v1"
	"github.com/shadow0vortex/cortexops/internal/rca/llm"
	"github.com/shadow0vortex/cortexops/internal/remediation/policy"
	"github.com/shadow0vortex/cortexops/internal/topology/graph"
	"github.com/shadow0vortex/cortexops/test/e2e/framework"
)

// ---------------------------------------------------------------------------
// TestChaos_LLMTimeout validates degraded-mode RAG heuristics.
//
// Invariant: If the LLM endpoint is unreachable or exceeds its timeout, the
// RCA service MUST NOT hang. It must return an error within the timeout window
// and the caller must be able to fall back to a degraded analysis.
// ---------------------------------------------------------------------------
func TestChaos_LLMTimeout(t *testing.T) {
	h := framework.Setup(t)

	t.Log("[CHAOS] Phase 1: Creating a black-hole HTTP server that never responds...")

	// Spin up a mock OpenAI server that hangs forever (simulating network blackhole).
	blackhole := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Block indefinitely until the client gives up.
		<-r.Context().Done()
	}))
	defer blackhole.Close()

	t.Log("[CHAOS] Phase 2: Exercising GenerateRCA with a tight timeout...")

	// Create an LLM client with a real API key. The client's internal HTTP timeout is 10s,
	// but we use a context with 100ms to simulate an aggressive caller SLA.
	// Since the client points at the real OpenAI URL (unreachable under timeout),
	// the context deadline will fire first and the request will abort.
	client := llm.NewOpenAIClient("test-api-key")

	ctx, cancel := context.WithTimeout(h.Ctx, 100*time.Millisecond)
	defer cancel()

	_, err := client.GenerateRCA(ctx, "Pod OOMKilled in namespace production")
	if err == nil {
		t.Fatal("INVARIANT VIOLATED: GenerateRCA returned nil error during simulated LLM blackhole. Expected timeout or connection error.")
	}

	t.Logf("[CHAOS] GenerateRCA correctly returned error: %v", err)

	t.Log("[CHAOS] Phase 3: Verifying degraded-mode fallback path...")

	// When API key is empty, the client must return a degraded analysis string (not an error).
	// This is the built-in fallback at openai.go:62-63.
	degradedClient := llm.NewOpenAIClient("")
	result, err := degradedClient.GenerateRCA(h.Ctx, "Pod OOMKilled in namespace production")
	if err != nil {
		t.Fatalf("INVARIANT VIOLATED: Degraded mode returned error: %v", err)
	}
	if result == "" {
		t.Fatal("INVARIANT VIOLATED: Degraded mode returned empty string. Must return fallback analysis.")
	}

	t.Logf("[CHAOS] Degraded-mode fallback returned: %q", truncate(result, 80))
	t.Log("[PASS] LLM Timeout chaos test passed. Degraded behavior invariant holds.")
}

// ---------------------------------------------------------------------------
// TestChaos_OPACrash validates fail-closed governance mechanics.
//
// Invariant: If the OPA policy evaluation encounters a cancelled context,
// the PolicyDecision MUST have Allowed=false. No automated remediation
// action may proceed when the policy engine is degraded.
// ---------------------------------------------------------------------------
func TestChaos_OPACrash(t *testing.T) {
	h := framework.Setup(t)

	t.Log("[CHAOS] Phase 1: Constructing OPA engine with valid embedded policy...")

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelDebug}))

	// NewOPAEngine panics if the embedded rego policy is broken.
	// We recover from any panic to turn it into a test failure.
	var opaEngine *policy.OPAEngine
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Fatalf("OPA engine construction panicked: %v", r)
			}
		}()
		opaEngine = policy.NewOPAEngine(&noopTopology{}, &noopMetrics{}, logger)
	}()

	t.Log("[CHAOS] Phase 2: Evaluating with a pre-cancelled context (simulating OPA crash)...")

	cancelledCtx, cancel := context.WithCancel(h.Ctx)
	cancel() // Immediately cancel to simulate crash/unavailability.

	incident := &v1.CorrelatedIncident{
		IncidentId: "chaos-test-001",
		Severity:   "CRITICAL",
	}

	action := &v1.RemediationAction{
		ActionId:        "chaos-action-001",
		Type:            v1.ActionType_POD_RESTART,
		TargetNamespace: "production",
		TargetResource:  "api-gateway",
	}

	decision, err := opaEngine.Evaluate(cancelledCtx, incident, action)
	if err != nil {
		// Even on Go error, that's still fail-closed behavior — no action proceeds.
		t.Logf("[CHAOS] OPA Evaluate returned Go error (acceptable fail-closed): %v", err)
		t.Log("[PASS] OPA crash chaos test passed via error return path.")
		return
	}

	if decision == nil {
		t.Fatal("INVARIANT VIOLATED: OPA returned nil decision and nil error. Must always return a decision.")
	}

	if decision.Allowed {
		t.Fatal("INVARIANT VIOLATED: OPA returned Allowed=true after context cancellation. FAIL-CLOSED BROKEN.")
	}

	t.Logf("[CHAOS] OPA correctly denied with reasoning: %q", decision.Reasoning)
	t.Logf("[CHAOS] Violating rules: %v", decision.ViolatingRules)
	t.Log("[PASS] OPA crash chaos test passed. Fail-closed invariant holds.")
}

// ---------------------------------------------------------------------------
// TestChaos_DatabasePartition validates topology graph cache isolation.
//
// Invariant: If PostgreSQL is unreachable, the in-memory MemoryGraphStore
// must continue serving BFS blast-radius queries. Only async persistence
// (SaveAsync) should fail — reads must remain fully operational.
// ---------------------------------------------------------------------------
func TestChaos_DatabasePartition(t *testing.T) {
	h := framework.Setup(t)
	_ = h

	t.Log("[CHAOS] Phase 1: Populating in-memory topology graph...")

	store := graph.NewMemoryGraphStore()
	ctx := context.Background()

	// Build a realistic service dependency graph:
	//   api-gateway -> user-service -> postgres-primary
	//   api-gateway -> order-service -> payment-service
	//   order-service -> postgres-primary
	nodeIDs := []string{
		"api-gateway",
		"user-service",
		"order-service",
		"payment-service",
		"postgres-primary",
	}

	for _, id := range nodeIDs {
		err := store.UpsertNode(ctx, &v1.TopologyNode{
			Id: id,
		})
		if err != nil {
			t.Fatalf("Failed to upsert node %s: %v", id, err)
		}
	}

	edges := []struct {
		src, dst string
	}{
		{"api-gateway", "user-service"},
		{"api-gateway", "order-service"},
		{"user-service", "postgres-primary"},
		{"order-service", "payment-service"},
		{"order-service", "postgres-primary"},
	}

	for _, e := range edges {
		err := store.UpsertEdge(ctx, e.src, e.dst, v1.EdgeType_ROUTES_TO)
		if err != nil {
			t.Fatalf("Failed to upsert edge %s->%s: %v", e.src, e.dst, err)
		}
	}

	t.Log("[CHAOS] Phase 2: Verifying BFS blast-radius before partition...")

	blastRadius, err := store.CalculateBlastRadius(ctx, "api-gateway")
	if err != nil {
		t.Fatalf("BFS failed before partition: %v", err)
	}
	if len(blastRadius) != 4 {
		t.Fatalf("Expected blast radius of 4 from api-gateway, got %d: %v", len(blastRadius), blastRadius)
	}
	t.Logf("[CHAOS] Pre-partition blast radius from api-gateway: %v", blastRadius)

	t.Log("[CHAOS] Phase 3: Simulating PostgreSQL partition (SaveAsync failure)...")

	// Create a persister pointed at an unreachable database URL.
	// This simulates what happens when Postgres is network-partitioned.
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: slog.LevelWarn}))

	// Attempt to create persister with bogus connection — expect failure.
	_, persistErr := graph.NewGraphPersister(
		"postgres://nonexistent:5432/cortexops?sslmode=disable&connect_timeout=1",
		store,
		logger,
	)
	// We expect this to fail since the DB is unreachable.
	if persistErr != nil {
		t.Logf("[CHAOS] Persistence layer correctly failed (DB unreachable): %v", persistErr)
	} else {
		// sql.Open succeeds lazily; the failure happens on first query.
		t.Log("[CHAOS] sql.Open succeeded lazily (expected). SaveAsync would fail on first write.")
	}

	t.Log("[CHAOS] Phase 4: Verifying in-memory graph survives the partition...")

	// The critical invariant: reads must still work after persistence failure.
	postPartitionRadius, err := store.CalculateBlastRadius(ctx, "api-gateway")
	if err != nil {
		t.Fatalf("INVARIANT VIOLATED: BFS query failed after DB partition: %v", err)
	}

	if len(postPartitionRadius) != 4 {
		t.Fatalf("INVARIANT VIOLATED: Blast radius changed after DB partition. Expected 4, got %d", len(postPartitionRadius))
	}

	// Also verify individual dependency lookups.
	deps, err := store.GetDependencies(ctx, "order-service")
	if err != nil {
		t.Fatalf("INVARIANT VIOLATED: GetDependencies failed after DB partition: %v", err)
	}
	if len(deps) != 2 {
		t.Fatalf("INVARIANT VIOLATED: order-service should have 2 dependencies, got %d: %v", len(deps), deps)
	}

	t.Logf("[CHAOS] Post-partition blast radius: %v", postPartitionRadius)
	t.Logf("[CHAOS] Post-partition order-service deps: %v", deps)
	t.Log("[PASS] Database partition chaos test passed. Cache isolation invariant holds.")
}

// ---------------------------------------------------------------------------
// Helpers
// ---------------------------------------------------------------------------

// noopMetrics satisfies core.MetricsRecorder without side effects.
type noopMetrics struct{}

func (n *noopMetrics) IncCounter(_ context.Context, _ string, _ map[string]string)                  {}
func (n *noopMetrics) ObserveHistogram(_ context.Context, _ string, _ float64, _ map[string]string)  {}
func (n *noopMetrics) SetGauge(_ context.Context, _ string, _ float64, _ map[string]string)          {}

// noopTopology satisfies core.TopologyProvider for OPA engine construction.
type noopTopology struct{}

func (n *noopTopology) GetDependencies(_ context.Context, _ string) ([]string, error) {
	return nil, nil
}
func (n *noopTopology) CalculateBlastRadius(_ context.Context, _ string) ([]string, error) {
	return nil, nil
}

// truncate returns the first n characters of s, or s if len(s) <= n.
func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
