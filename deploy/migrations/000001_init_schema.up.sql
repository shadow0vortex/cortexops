CREATE TABLE IF NOT EXISTS topology_snapshots (
    id SERIAL PRIMARY KEY,
    version BIGINT NOT NULL,
    checksum VARCHAR(64) NOT NULL,
    data JSONB NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_topology_snapshots_version ON topology_snapshots(version);
