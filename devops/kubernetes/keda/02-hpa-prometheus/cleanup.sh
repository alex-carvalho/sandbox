#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "Deleting namespace hpa-poc..."
kubectl delete -f "${SCRIPT_DIR}/manifests/hpa.yaml" --ignore-not-found
kubectl delete -f "${SCRIPT_DIR}/manifests/consumer.yaml" --ignore-not-found
kubectl delete -f "${SCRIPT_DIR}/manifests/namespace.yaml" --ignore-not-found
echo "HPA workload cleaned up."
