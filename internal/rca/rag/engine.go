package rag

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	correlationv1 "github.com/cortexops/cortexops/api/v1"
	rcav1 "github.com/cortexops/cortexops/api/v1"
	rcactx "github.com/cortexops/cortexops/internal/rca/context"
	"github.com/cortexops/cortexops/pkg/core"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Engine is the RAG orchestrator for incident RCA.
type Engine struct {
	builder   *rcactx.Builder
	embedder  core.EmbeddingClient
	llm       core.LLMClient
	vectorDB  core.VectorStore
	publisher core.Publisher
	metrics   core.MetricsRecorder
	logger    *slog.Logger
}

func NewEngine(embedder core.EmbeddingClient, llm core.LLMClient, vectorDB core.VectorStore, pub core.Publisher, metrics core.MetricsRecorder, logger *slog.Logger) *Engine {
	return &Engine{
		builder:   rcactx.NewBuilder(2048), // Max context tokens
		embedder:  embedder,
		llm:       llm,
		vectorDB:  vectorDB,
		publisher: pub,
		metrics:   metrics,
		logger:    logger,
	}
}

// GenerateRCA processes a CorrelatedIncident and deterministically builds an AI-assisted RCAReport.
func (e *Engine) GenerateRCA(ctx context.Context, incident *correlationv1.CorrelatedIncident) error {
	start := time.Now()
	defer func() {
		e.metrics.ObserveHistogram(ctx, "cortexops_ai_inference_latency_seconds", time.Since(start).Seconds(), nil)
	}()

	// 1. Build context from raw incident
	incidentStr, evidenceIDs, truncated := e.builder.Build(incident)

	// 2. Fetch Embeddings and Historical Context
	// We wrap this in timeout and fallback to Degraded Mode on failure to guarantee delivery.
	var historical []*rcav1.HistoricalSimilarity
	embedCtx, embedCancel := context.WithTimeout(ctx, 3*time.Second)
	defer embedCancel()

	embedding, err := e.embedder.GenerateEmbeddings(embedCtx, incidentStr)
	if err == nil {
		historical, _ = e.vectorDB.SearchSimilar(embedCtx, embedding, 3, "")
	} else {
		e.logger.Warn("Embedding service unreachable, degrading RCA", "error", err)
	}

	// 3. Prompt Assembly
	prompt := e.assemblePrompt(incidentStr, historical)

	// 4. LLM Generation
	llmCtx, llmCancel := context.WithTimeout(ctx, 15*time.Second)
	defer llmCancel()

	analysis, err := e.llm.GenerateRCA(llmCtx, prompt)
	isDegraded := false
	if err != nil {
		e.logger.Error("LLM generation failed, publishing degraded RCA", "error", err)
		isDegraded = true
		analysis = "AI Analysis unavailable due to inference failure. Rely on attached telemetry evidence."
	}

	// 5. Publish Output (Immutable, Traceable)
	report := &rcav1.RCAReport{
		RcaId:       uuid.New().String(),
		IncidentId:  incident.IncidentId,
		GeneratedAt: timestamppb.Now(),
		Analysis:    analysis,
		GroundedEvidence: &rcav1.RCAEvidence{
			TelemetryEventIds: evidenceIDs,
			WasTruncated:      truncated,
		},
		HistoricalContext: historical,
		IsDegraded:        isDegraded,
	}

	e.logger.Info("RCA generation complete", "incident_id", incident.IncidentId, "degraded", isDegraded)
	
	// Payload marshaling and publishing would happen here.
	return nil // err omitted for brevity
}

func (e *Engine) assemblePrompt(incident string, history []*rcav1.HistoricalSimilarity) string {
	prompt := "You are CortexOps, an advisory AIOps system. Analyze the following Kubernetes incident deterministic telemetry.\n\n"
	prompt += incident + "\n\n"

	if len(history) > 0 {
		prompt += "--- HISTORICAL CONTEXT ---\n"
		for _, h := range history {
			prompt += fmt.Sprintf("Similar Incident Resolved By: %s\n", h.ResolutionSummary)
		}
	}

	prompt += "\nINSTRUCTIONS:\n"
	prompt += "1. Summarize the root cause deterministically based ONLY on the evidence provided.\n"
	prompt += "2. Suggest operational actions (Read-Only) to verify the hypothesis.\n"
	prompt += "3. DO NOT output code or kubectl commands."
	return prompt
}
