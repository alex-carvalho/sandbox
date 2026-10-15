#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=========================================================="
echo "  Installing Common Infrastructure"
echo "  - Namespaces: common, keda, monitoring"
echo "  - Redis + redis_exporter (keda-queue & hpa-prometheus-queue)"
echo "  - KEDA Operator & Metrics Server"
echo "  - Prometheus Server & Prometheus Adapter"
echo "=========================================================="

echo "Ensuring kubectl context is 'kind-keda-vs-hpa'..."
kubectl config use-context "kind-keda-vs-hpa"

# 1. Apply namespaces
echo "1. Applying namespaces..."
kubectl apply -f "${SCRIPT_DIR}/manifests/namespace.yaml"

# 2. Deploy Redis with redis_exporter
echo "2. Deploying Redis + Exporter in 'common' namespace..."
kubectl apply -f "${SCRIPT_DIR}/manifests/redis-with-exporter.yaml"

echo "Waiting for Redis to be ready..."
kubectl rollout status deployment/redis -n common --timeout=120s

# Initialize queues with an item and drain to ensure keys are recognized
echo "Initializing Redis keys in background..."
kubectl exec -n common deployment/redis -c redis -- redis-cli rpush keda-queue init >/dev/null 2>&1 || true
kubectl exec -n common deployment/redis -c redis -- redis-cli lpop keda-queue >/dev/null 2>&1 || true
kubectl exec -n common deployment/redis -c redis -- redis-cli rpush hpa-prometheus-queue init >/dev/null 2>&1 || true
kubectl exec -n common deployment/redis -c redis -- redis-cli lpop hpa-prometheus-queue >/dev/null 2>&1 || true

# 3. Add and update Helm repositories
echo "3. Updating Helm repositories..."
helm repo add kedacore https://kedacore.github.io/charts --force-update
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts --force-update
helm repo update

# 4. Install KEDA
echo "4. Installing KEDA in 'keda' namespace..."
helm upgrade --install keda kedacore/keda \
  --namespace keda \
  --set watchNamespace="" \
  --wait \
  --timeout 180s

# 5. Install Prometheus
echo "5. Installing Prometheus in 'monitoring' namespace..."
helm upgrade --install prometheus prometheus-community/prometheus \
  --namespace monitoring \
  -f "${SCRIPT_DIR}/manifests/prometheus-values.yaml" \
  --wait \
  --timeout 180s

# 6. Install Prometheus Adapter
echo "6. Installing Prometheus Adapter in 'monitoring' namespace..."
helm upgrade --install prometheus-adapter prometheus-community/prometheus-adapter \
  --namespace monitoring \
  -f "${SCRIPT_DIR}/manifests/prometheus-adapter-values.yaml" \
  --wait \
  --timeout 180s

echo ""
echo "=========================================================="
echo "  Common Infrastructure Successfully Deployed!"
echo "=========================================================="
kubectl get pods -n common
kubectl get pods -n keda
kubectl get pods -n monitoring
