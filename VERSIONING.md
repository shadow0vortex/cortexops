# Versioning Policy

CortexOps strictly follows Semantic Versioning 2.0.0.

## Stability Guarantees
- **Protobuf Contracts**: Once a version reaches v1.0.0, any breaking changes to Protobuf definitions (`api/v1/`) will trigger a major version bump.
- **Workflow State**: We aim for backward compatibility in Temporal workflow history, but major version changes may require a clean start of remediation workflows.

## Deprecation Policy
Features slated for removal will be marked as `DEPRECATED` in the documentation and logs for at least one minor release before being removed in a major release.
