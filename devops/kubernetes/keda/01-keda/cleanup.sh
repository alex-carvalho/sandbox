#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Deleting namespace keda-poc..."
kubectl delete -f "${SCRIPT_DIR}/manifests/scaledobject.yaml" --ignore-not-found
kubectl delete -f "${SCRIPT_DIR}/manifests/consumer.yaml" --ignore-not-found
kubectl delete -f "${SCRIPT_DIR}/manifests/namespace.yaml" --ignore-not-found
echo "KEDA workload cleaned up."
