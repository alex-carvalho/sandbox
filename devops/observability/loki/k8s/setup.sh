#!/bin/bash
set -euo pipefail

CLUSTER_NAME="loki-obs-poc"

if ! command -v kind &>/dev/null; then
  echo "kind is not installed"
  exit 1
fi

if ! command -v kubectl &>/dev/null; then
  echo "kubectl is not installed"
  exit 1
fi

if ! command -v helm &>/dev/null; then
  echo "helm is not installed"
  exit 1
fi

if kind get clusters | grep -q "^${CLUSTER_NAME}$"; then
  NODE_COUNT=$(kubectl get nodes --no-headers 2>/dev/null | wc -l | tr -d ' ' || echo "0")
  if [ "$NODE_COUNT" -eq 3 ]; then
    echo "Kind cluster ${CLUSTER_NAME} already exists with 3 nodes. Skipping creation."
  else
    echo "Recreating Kind cluster ${CLUSTER_NAME} to apply 3-node configuration..."
    kind delete cluster --name "$CLUSTER_NAME"
    kind create cluster --name "$CLUSTER_NAME" --config kind-config.yaml
  fi
else
  kind create cluster --name "$CLUSTER_NAME" --config kind-config.yaml
fi

helm repo add grafana https://grafana.github.io/helm-charts --force-update
helm repo add open-telemetry https://open-telemetry.github.io/opentelemetry-helm-charts --force-update
helm repo update

kubectl create namespace observability --dry-run=client -o yaml | kubectl apply -f -

helm upgrade --install loki grafana/loki \
  --namespace observability \
  -f helm/loki-values.yaml

helm upgrade --install grafana grafana/grafana \
  --namespace observability \
  -f helm/grafana-values.yaml

kubectl apply -f manifests/grafana-dashboard.yaml

helm upgrade --install otel-collector open-telemetry/opentelemetry-collector \
  --namespace observability \
  -f helm/otel-collector-values.yaml

kubectl apply -f manifests/demo-app.yaml

echo "Waiting for Loki SimpleScalable components to be ready..."
kubectl rollout status deployment/loki-gateway -n observability --timeout=300s
kubectl rollout status statefulset/loki-write -n observability --timeout=300s
kubectl rollout status deployment/loki-read -n observability --timeout=300s
kubectl rollout status statefulset/loki-backend -n observability --timeout=300s

echo "Waiting for Grafana to be ready..."
kubectl rollout status deployment/grafana -n observability --timeout=300s

echo "Waiting for OpenTelemetry Collector to be ready..."
kubectl rollout status daemonset/otel-collector-opentelemetry-collector-agent -n observability --timeout=300s

echo "Waiting for log generator app to be ready..."
kubectl rollout status deployment/log-generator -n demo-apps --timeout=300s

echo "All components are deployed and running successfully!"
echo "To access Grafana dashboard, run:"
echo "  kubectl port-forward svc/grafana 3000:80 -n observability"
echo "Then visit: http://localhost:3000"
echo "Login: admin / admin"
