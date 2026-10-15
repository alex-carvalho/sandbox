#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"
CLUSTER_NAME="keda-vs-hpa"

echo "=========================================================="
echo "  Setting up Kind Cluster: ${CLUSTER_NAME}"
echo "=========================================================="

# 1. Pre-flight checks
if ! command -v docker &>/dev/null; then
  echo "Error: docker is not installed or not in PATH."
  exit 1
fi

if ! docker info &>/dev/null; then
  echo "Docker daemon is not responding."
  if command -v colima &>/dev/null; then
    echo "Attempting to start colima..."
    colima start
  else
    echo "Please start your Docker daemon and retry."
    exit 1
  fi
fi

if ! command -v kind &>/dev/null; then
  echo "Error: kind CLI is not installed."
  exit 1
fi

if ! command -v kubectl &>/dev/null; then
  echo "Error: kubectl CLI is not installed."
  exit 1
fi

if ! command -v helm &>/dev/null; then
  echo "Error: helm CLI is not installed."
  exit 1
fi

# 2. Create Kind cluster if it doesn't already exist
if kind get clusters | grep -q "^${CLUSTER_NAME}$"; then
  echo "Kind cluster '${CLUSTER_NAME}' already exists. Switching context..."
  kubectl config use-context "kind-${CLUSTER_NAME}"
else
  echo "Creating Kind cluster '${CLUSTER_NAME}'..."
  kind create cluster --config "${SCRIPT_DIR}/kind-config.yaml"
fi

# Verify cluster nodes
echo "Waiting for cluster nodes to be Ready..."
kubectl wait --for=condition=Ready nodes --all --timeout=60s

# 3. Bootstrap Common Infrastructure
echo ""
echo "Bootstrapping common infrastructure (Redis, Exporter, KEDA, Prometheus)..."
"${ROOT_DIR}/common/install-common.sh"

echo ""
echo "=========================================================="
echo "  Setup Complete! Cluster is ready for both PoC workloads."
echo "=========================================================="
echo "Next steps:"
echo "  - To deploy & test KEDA:            cd 01-keda && ./deploy.sh && ./test.sh"
echo "  - To deploy & test HPA+Prometheus:  cd 02-hpa-prometheus && ./deploy.sh && ./test.sh"
echo "  - Or run load generator directly:  ./common/load-generator.sh status"
