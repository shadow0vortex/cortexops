# Security Policy

## Supported Versions

We provide security updates for the following versions of CortexOps:

| Version | Supported          |
| ------- | ------------------ |
| v1.0.x  | :white_check_mark: |
| < v1.0  | :x:                |

## Reporting a Vulnerability

We take the security of CortexOps seriously. If you believe you have found a security vulnerability, please report it to us by following these steps:

1.  **Do not open a public issue.**
2.  Email your report to `security@cortexops.io` (Note: This is a placeholder for demo purposes).
3.  Include as much detail as possible, including steps to reproduce and a Proof of Concept (PoC) if available.

We will acknowledge your report within 48 hours and provide a timeline for resolution.

## Operational Safety Note
CortexOps modifies Kubernetes cluster state. Users are responsible for ensuring that OPA policies and Temporal workflow timeouts are correctly configured for their specific environment. Always test remediation logic in a non-production environment first.
