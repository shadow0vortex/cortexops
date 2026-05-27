#!/bin/bash
set -e

echo "Starting CortexOps Repository Hygiene Verification..."

# 1. Verify clean protobuf generation (no uncommitted diffs in api/v1)
if ! git diff --quiet -- api/v1/; then
  echo "ERROR: Uncommitted protobuf generation changes detected. Run 'make proto' and commit the results."
  git diff --name-only api/v1/
  exit 1
fi

# 2. Prevent stale mocks from reappearing in runtime code (excluding tests)
if grep -r -i "mock" --include="*.go" --exclude="*_test.go" cmd/ internal/ pkg/; then
  echo "ERROR: Found mock implementations in runtime code. Mocks are only allowed in tests."
  exit 1
fi

# 3. Ensure build compiles cleanly
echo "Compiling all modules..."
if ! go build ./...; then
  echo "ERROR: Build failed. Compilation errors present."
  exit 1
fi

# 4. Ensure tests compile and pass
echo "Running unit tests..."
if ! go test ./...; then
  echo "ERROR: Tests failed."
  exit 1
fi

echo "Repository Hygiene Verification Passed!"
