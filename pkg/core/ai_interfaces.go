package core

import (
	"context"

	rcav1 "github.com/cortexops/cortexops/api/v1"
)

// VectorStore abstracts vector database operations (e.g., Qdrant) for historical memory.
type VectorStore interface {
	UpsertIncident(ctx context.Context, incidentID string, embedding []float32, metadata map[string]string) error
	SearchSimilar(ctx context.Context, embedding []float32, limit int, namespace string) ([]*rcav1.HistoricalSimilarity, error)
}

// EmbeddingClient abstracts models for generating text embeddings.
type EmbeddingClient interface {
	GenerateEmbeddings(ctx context.Context, text string) ([]float32, error)
}

// LLMClient abstracts Language Model generation.
type LLMClient interface {
	GenerateRCA(ctx context.Context, prompt string) (string, error)
}
