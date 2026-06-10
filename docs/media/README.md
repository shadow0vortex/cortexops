# Screenshot & Media Assets

This directory should be populated with professional screenshots for the GitHub repository and presentation decks.

## Recommended Assets

### 1. `grafana_overview.png`
- **Source**: `http://localhost:3000`
- **View**: `CortexOps Demo Overview` dashboard.
- **Goal**: Show active incidents, ingestion rates, and remediation status.

### 2. `temporal_workflow_history.png`
- **Source**: `http://localhost:8233`
- **View**: A specific execution of `RemediationWorkflow`.
- **Goal**: Show the durability and step-by-step audit trail (Policy -> DryRun -> Execute).

### 3. `topology_graph_vis.png`
- **Source**: `make diagnostics` (transformed into a visualization)
- **Goal**: Show the directed graph relationships between cluster nodes.

### 4. `remediation_audit.png`
- **Source**: `curl -s http://localhost:9091/debug/incidents/active`
- **Goal**: Show structured JSON evidence associated with a correlated incident.
