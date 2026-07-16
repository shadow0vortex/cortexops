-- =============================================================================
-- Rollback: Per-Service PostgreSQL Schemas (TD-039)
-- Moves tables back to public and drops per-service schemas.
-- =============================================================================

-- Move topology_snapshots back to public schema
ALTER TABLE IF EXISTS topology.topology_snapshots SET SCHEMA public;

-- Drop per-service schemas
DROP SCHEMA IF EXISTS remediation CASCADE;
DROP SCHEMA IF EXISTS rca CASCADE;
DROP SCHEMA IF EXISTS correlator CASCADE;
DROP SCHEMA IF EXISTS topology CASCADE;
