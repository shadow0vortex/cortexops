# Security Pipeline

This document outlines the automated security scanning and vulnerability management processes implemented in CortexOps via `security.yml`.

## Philosophy

The objective of the security pipeline is to prevent **newly introduced** regressions without indefinitely blocking the repository because of legacy or accepted risk vulnerabilities. All security scans are designed to be informative and non-blocking, prioritizing continuous delivery while ensuring security visibility.

## Tools Integrated

1. **`govulncheck`**
   - **Target**: Go dependencies and standard library.
   - **Behavior**: Uses call-graph analysis to identify if vulnerable functions in dependencies are actually executed in the CortexOps codebase.
   - **Configuration**: Runs universally across `./...`.

2. **`gosec` (Go Security)**
   - **Target**: Go source code (SAST).
   - **Behavior**: Scans for hardcoded credentials, unsafe `unsafe` pointer usage, SQL injection, and weak random number generators.
   - **Configuration**: Uses `-no-fail` and exports to SARIF format for GitHub Advanced Security integration.

3. **`Trivy`**
   - **Target**: Filesystem scanning.
   - **Behavior**: Identifies outdated dependencies in non-Go contexts and provides high-level OS vulnerability context for the Dockerfiles.
   - **Configuration**: Triggers on `CRITICAL,HIGH` severities and outputs SARIF files.

4. **`Gitleaks`**
   - **Target**: Commit history and diffs.
   - **Behavior**: Prevents accidental committal of API keys, AWS credentials, JWT tokens, and private keys.

## Artifact Handling

If GitHub Advanced Security (GHAS) is enabled on the repository:
- SARIF files are uploaded directly to GitHub via `github/codeql-action/upload-sarif`.
- Vulnerabilities are visible in the "Security" tab.

If GHAS is unavailable:
- The pipeline securely uploads the raw `gosec.sarif` and `trivy-results.sarif` as standard workflow artifacts (`actions/upload-artifact`) retained for 7 days. These can be downloaded and analyzed manually.
