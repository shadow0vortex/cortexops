-- =============================================================================
-- Migration: Per-Service PostgreSQL Schemas (TD-039)
-- Creates isolated schemas for each CortexOps service within the shared DB.
-- =============================================================================

-- Create per-service schemas
CREATE SCHEMA IF NOT EXISTS topology;
CREATE SCHEMA IF NOT EXISTS correlator;
CREATE SCHEMA IF NOT EXISTS rca;
CREATE SCHEMA IF NOT EXISTS remediation;

-- Move existing topology_snapshots table into the topology schema
ALTER TABLE IF EXISTS public.topology_snapshots SET SCHEMA topology;
