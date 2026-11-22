#!/usr/bin/env bash
set -euo pipefail

echo "Checking Temporal UI..."
MAX_RETRIES=15
RETRY_COUNT=0
until curl --fail --silent --show-error http://localhost:8233/ >/dev/null || [ "$RETRY_COUNT" -ge "$MAX_RETRIES" ]; do
    echo "Waiting for Temporal UI..."
    sleep 5
    RETRY_COUNT=$((RETRY_COUNT + 1))
done

if [ "$RETRY_COUNT" -ge "$MAX_RETRIES" ]; then
    echo "ERROR: Temporal UI is not responding on port 8233" >&2
    exit 1
fi

echo "Checking Temporal default namespace..."
docker run --rm --network cortexops-dev \
    -e TEMPORAL_ADDRESS=temporal:7233 \
    --entrypoint temporal temporalio/admin-tools:1.23.0 \
    operator namespace describe default >/dev/null

echo "Running the live failure-to-workflow validation..."
make validate-pipeline
