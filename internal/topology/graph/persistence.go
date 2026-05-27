package graph

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	_ "github.com/lib/pq"
	topologyv1 "github.com/shadow0vortex/cortexops/api/v1"
)

type SnapshotData struct {
	Nodes map[string]*topologyv1.TopologyNode   `json:"nodes"`
	Edges map[string][]*topologyv1.TopologyEdge `json:"edges"`
}

type GraphPersister struct {
	db     *sql.DB
	logger *slog.Logger
	store  *MemoryGraphStore
}

func NewGraphPersister(dbURL string, store *MemoryGraphStore, logger *slog.Logger) (*GraphPersister, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open postgres: %w", err)
	}

	// Ensure table exists
	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS topology_snapshots (
		id SERIAL PRIMARY KEY,
		version INT NOT NULL,
		checksum VARCHAR(64) NOT NULL,
		data JSONB NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to create snapshot table: %w", err)
	}

	return &GraphPersister{
		db:     db,
		store:  store,
		logger: logger,
	}, nil
}

// Restore loads the latest snapshot into memory and verifies checksum.
func (p *GraphPersister) Restore(ctx context.Context) error {
	var dataBytes []byte
	var checksum string
	var version int

	err := p.db.QueryRowContext(ctx, `
		SELECT version, checksum, data 
		FROM topology_snapshots 
		ORDER BY id DESC LIMIT 1
	`).Scan(&version, &checksum, &dataBytes)

	if err == sql.ErrNoRows {
		p.logger.Info("No existing topology snapshot found, starting fresh")
		return nil
	} else if err != nil {
		return fmt.Errorf("failed to load snapshot: %w", err)
	}

	// Verify checksum
	hash := fmt.Sprintf("%x", sha256.Sum256(dataBytes))
	if hash != checksum {
		return fmt.Errorf("snapshot corruption detected! Expected checksum %s but got %s", checksum, hash)
	}

	var data SnapshotData
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		return fmt.Errorf("failed to unmarshal snapshot data: %w", err)
	}

	// Load into memory store
	p.store.mu.Lock()
	defer p.store.mu.Unlock()

	p.store.nodes = make(map[string]*topologyv1.TopologyNode)
	p.store.edges = make(map[string][]*topologyv1.TopologyEdge)

	for k, v := range data.Nodes {
		p.store.nodes[k] = v
	}
	for k, v := range data.Edges {
		p.store.edges[k] = v
	}

	p.logger.Info("Successfully restored topology snapshot", "version", version, "nodes_count", len(data.Nodes))
	return nil
}

// SaveAsync takes a snapshot and saves it to postgres.
func (p *GraphPersister) SaveAsync(ctx context.Context) error {
	dataBytes, err := p.store.Snapshot(ctx)
	if err != nil {
		return fmt.Errorf("failed to snapshot memory store: %w", err)
	}

	hash := fmt.Sprintf("%x", sha256.Sum256(dataBytes))
	version := time.Now().Unix()

	_, err = p.db.ExecContext(ctx, `
		INSERT INTO topology_snapshots (version, checksum, data) 
		VALUES ($1, $2, $3)
	`, version, hash, string(dataBytes))

	if err != nil {
		return fmt.Errorf("failed to insert snapshot: %w", err)
	}

	p.logger.Debug("Saved topology snapshot", "version", version, "checksum", hash)
	return nil
}

func (p *GraphPersister) Close() error {
	return p.db.Close()
}
