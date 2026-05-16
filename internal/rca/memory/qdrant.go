package memory

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	rcav1 "github.com/cortexops/cortexops/api/v1"
)

// QdrantClient implements core.VectorStore using the Qdrant REST API.
type QdrantClient struct {
	url        string
	collection string
	httpClient *http.Client
}

func NewQdrantClient(url, collection string) *QdrantClient {
	return &QdrantClient{
		url:        url,
		collection: collection,
		httpClient: &http.Client{},
	}
}

// UpsertIncident saves the incident embedding for future RCA correlations.
func (q *QdrantClient) UpsertIncident(ctx context.Context, incidentID string, embedding []float32, metadata map[string]string) error {
	payload := map[string]interface{}{
		"points": []map[string]interface{}{
			{
				"id":      incidentID, // Assumes UUID format
				"vector":  embedding,
				"payload": metadata,
			},
		},
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "PUT", fmt.Sprintf("%s/collections/%s/points", q.url, q.collection), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := q.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("qdrant connection failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("qdrant rejected points: %s", string(respBody))
	}
	return nil
}

// SearchSimilar fetches past incidents to enrich the RAG prompt.
func (q *QdrantClient) SearchSimilar(ctx context.Context, embedding []float32, limit int, namespace string) ([]*rcav1.HistoricalSimilarity, error) {
	// A robust implementation would use query filters on the namespace payload.
	payload := map[string]interface{}{
		"vector": embedding,
		"limit":  limit,
		"with_payload": true,
	}

	body, _ := json.Marshal(payload)
	req, _ := http.NewRequestWithContext(ctx, "POST", fmt.Sprintf("%s/collections/%s/points/search", q.url, q.collection), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := q.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("qdrant search failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("qdrant search rejected")
	}

	// Parse response payload and construct HistoricalSimilarity list...
	// Stubbing parsing for scaffolding purposes.
	
	return []*rcav1.HistoricalSimilarity{}, nil
}
