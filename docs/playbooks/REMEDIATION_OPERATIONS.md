# Playbook: Remediation Operations

Guidelines for managing autonomous infrastructure mutations.

## 1. Safety Procedures
- **Dry-Run Validation**: Always observe the `DRY_RUNNING` state in the Temporal UI before the `EXECUTING` phase.
- **Rollback Verification**: If a remediation results in a `ROLLBACK`, do not re-run the same action without investigating the stabilization failure first.

## 2. Managing Policy Denials
If an action is blocked by OPA:
1. Identify the violating rule in the `AuditRecord`.
2. To override, update the OPA policy in `internal/remediation/policy/` or adjust the incident severity.

## 3. Human-in-the-Loop
For high-risk actions (Risk Score > 0.7), monitor the approval queue. Remediations will time out after 1 hour if no approval is received.
