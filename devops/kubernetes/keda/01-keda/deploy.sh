#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================================="
echo "  Deploying KEDA Workload (Namespace: keda-poc)"
echo "=========================================================="

kubectl apply -f "${SCRIPT_DIR}/manifests/namespace.yaml"
kubectl apply -f "${SCRIPT_DIR}/manifests/consumer.yaml"
kubectl apply -f "${SCRIPT_DIR}/manifests/scaledobject.yaml"

echo ""
echo "KEDA workload deployed successfully!"
echo "Current deployment status (notice replicas: 0):"
kubectl get deployment keda-consumer -n keda-poc
kubectl get scaledobject keda-queue-scaler -n keda-poc
