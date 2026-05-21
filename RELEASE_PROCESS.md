# Release Process

This document describes how to release a new version of CortexOps.

## 1. Versioning Strategy
We use [Semantic Versioning](https://semver.org/).
- **MAJOR**: Breaking architectural changes or significant protobuf updates.
- **MINOR**: New remediation activities, correlation heuristics, or observability features.
- **PATCH**: Bug fixes, linting improvements, and documentation updates.

## 2. Release Steps

1.  **Draft Changelog**: Update `CHANGELOG.md` with all changes since the last release.
2.  **Tag the Release**:
    ```bash
    git tag -a v1.0.0 -m "Release v1.0.0"
    git push origin v1.0.0
    ```
3.  **Build Container Images**:
    GitHub Actions will automatically trigger a build and push to Docker Hub (or GHCR).
4.  **Publish GitHub Release**: Use the GitHub UI to create a new release from the tag, attaching the generated SBOM.

## 3. Post-Release
- Verify the new Helm chart version is deployable.
- Announce the release in the discussions or community channel.
