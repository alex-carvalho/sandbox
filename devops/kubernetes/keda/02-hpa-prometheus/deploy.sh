#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================================="
echo "  Deploying HPA + Prometheus Workload (Namespace: hpa-poc)"
echo "=========================================================="

kubectl apply -f "${SCRIPT_DIR}/manifests/namespace.yaml"
kubectl apply -f "${SCRIPT_DIR}/manifests/consumer.yaml"
kubectl apply -f "${SCRIPT_DIR}/manifests/hpa.yaml"

echo ""
echo "HPA workload deployed successfully!"
echo "Current deployment status (notice baseline replicas: 1):"
kubectl get deployment hpa-consumer -n hpa-poc
kubectl get hpa hpa-prometheus-consumer -n hpa-poc
