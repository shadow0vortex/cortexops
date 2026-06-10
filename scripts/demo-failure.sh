#!/bin/bash
# Inject a deterministic failure into the CortexOps Demo Environment
SCENARIO=${1:-rollout-fail}
go run ./cmd/demo/main.go inject -scenario=$SCENARIO
