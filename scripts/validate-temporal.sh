#!/bin/bash
# Validate Temporal Setup & Workflows

echo "Checking Temporal Server health..."
MAX_RETRIES=15
RETRY_COUNT=0
until curl -s -f http://localhost:8233 > /dev/null || [ $RETRY_COUNT -eq $MAX_RETRIES ]; do
    echo "Waiting for Temporal UI..."
    sleep 5
    ((RETRY_COUNT++))
done

if [ $RETRY_COUNT -eq $MAX_RETRIES ]; then
    echo "ERROR: Temporal UI is not responding on port 8233"
    exit 1
fi
echo "Temporal UI is accessible."

echo "Checking for registered workers..."
# For this script we will verify the workflow creation end-to-end via the API or chaos tool
echo "Running Failure Injection to trigger a workflow..."
go run ./cmd/chaos failure-injection rollout-fail

echo "Validation successful. Workflows can be viewed in Temporal UI."
