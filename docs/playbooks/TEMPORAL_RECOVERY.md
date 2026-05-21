# Playbook: Temporal Recovery

Temporal is the brain of our remediation orchestration. If it fails, all autonomous actions stall.

## 1. Database Connectivity
**Problem**: Temporal cannot connect to Postgres.
**Resolution**:
1. Check Postgres logs: `docker compose logs postgres`.
2. Ensure migrations have completed: `docker compose logs temporal`.

## 2. Worker Death
**Problem**: Remediation worker crashed.
**Resolution**: 
1. Restart the service: `make dev-up`.
2. Temporal will automatically deliver the current task to the new worker. Execution is guaranteed to resume from the last successful state.

## 3. Workflow Termination
If a workflow is truly stuck due to a logic bug:
1. Terminate it via the UI (`localhost:8233`).
2. Fix the bug, rebuild images, and restart the system.
