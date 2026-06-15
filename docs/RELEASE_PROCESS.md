# Release Process

This document describes the automated release pipeline for the CortexOps platform, governed by `.github/workflows/release.yml`.

## Trigger Mechanism

Releases are fully automated and driven by Git tags. 
When a new tag matching the pattern `v*` (e.g., `v1.0.0`, `v0.2.0-beta`) is pushed to the repository, the release workflow initiates.

## Automation Pipeline

1. **Environment Setup**: Provisions Ubuntu runner and caches the Go 1.22 environment.
2. **Cross-Compilation**: Iterates through every CortexOps microservice (`collector`, `correlator`, `topology`, `rca`, `remediation`).
3. **Architecture Targets**: Builds binaries for `linux/amd64` and `linux/arm64`.
4. **Checksum Generation**: Executes `sha256sum` against all compiled binaries to generate a `checksums.txt` file for verifiable integrity.
5. **GitHub Release Creation**: 
   - Uses `softprops/action-gh-release` to create a formal release.
   - Attaches all compiled binaries and the `checksums.txt`.
   - Automatically generates a changelog (`generate_release_notes: true`) based on merged Pull Requests since the last tag.

## Releasing a New Version

1. Ensure the `main` branch is stable and all Quality Gates have passed.
2. Execute the following commands:
   ```bash
   git checkout main
   git pull origin main
   git tag v1.0.0
   git push origin v1.0.0
   ```
3. The GitHub Actions release pipeline will build the assets and publish them automatically within a few minutes.
