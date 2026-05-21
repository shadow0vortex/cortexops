package topology

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// HTTPTopologyClient implements core.TopologyProvider by querying a remote diagnostics API.
type HTTPTopologyClient struct {
	baseURL string
	client  *http.Client
}

func NewHTTPClient(baseURL string) *HTTPTopologyClient {
	return &HTTPTopologyClient{
		baseURL: baseURL,
		client:  &http.Client{},
	}
}

func (c *HTTPTopologyClient) GetDependencies(ctx context.Context, nodeID string) ([]string, error) {
	u := fmt.Sprintf("%s/debug/graph/node?id=%s", c.baseURL, url.QueryEscape(nodeID))
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("topology API returned status %d", resp.StatusCode)
	}

	var result struct {
		Dependencies []string `json:"dependencies"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Dependencies, nil
}

func (c *HTTPTopologyClient) CalculateBlastRadius(ctx context.Context, nodeID string) ([]string, error) {
	u := fmt.Sprintf("%s/debug/graph/blast-radius?id=%s", c.baseURL, url.QueryEscape(nodeID))
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("topology API returned status %d", resp.StatusCode)
	}

	var result struct {
		ImpactedIDs []string `json:"impacted_ids"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.ImpactedIDs, nil
}
