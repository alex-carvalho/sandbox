#!/bin/bash

set -e

# Get the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
K8S_DIR="${SCRIPT_DIR}/k8s"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Helper functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

wait_for_resource() {
    local resource_type=$1
    local resource_name=$2
    local namespace=$3
    local timeout=${4:-300}
    
    log_info "Waiting for $resource_type/$resource_name in namespace $namespace..."
    
    kubectl wait --for=condition=ready $resource_type/$resource_name -n $namespace --timeout=${timeout}s || {
        log_warn "Timeout waiting for $resource_type/$resource_name"
        return 1
    }
}

# Step 1: Create KinD cluster
log_info "Creating KinD cluster..."
kind create cluster --name elastic-kind-cluster || {
    log_warn "Cluster already exists, skipping creation"
}
sleep 5

# Step 2: Create namespaces
log_info "Creating namespaces..."
kubectl create namespace elastic --dry-run=client -o yaml | kubectl apply -f -
kubectl create namespace apps --dry-run=client -o yaml | kubectl apply -f -

# Step 3: Add Elastic Helm repository and install ECK operator
log_info "Adding Elastic Helm repository..."
helm repo add elastic https://helm.elastic.co
helm repo update

log_info "Installing ECK operator..."
if helm list -n elastic | grep -q eck-operator; then
    log_info "ECK operator already installed, upgrading..."
    helm upgrade eck-operator elastic/eck-operator \
        --namespace elastic \
        --version 3.2.0 \
        --values - <<EOF
operator:
  watchNamespaces: []
EOF
else
    helm install eck-operator elastic/eck-operator \
        --namespace elastic \
        --version 3.2.0 \
        --values - <<EOF
operator:
  watchNamespaces: []
EOF
fi

# Step 4: Wait for ECK operator to be ready
log_info "Waiting for ECK operator StatefulSet to be ready (this may take a few minutes)..."
kubectl rollout status statefulset/elastic-operator -n elastic --timeout=600s || {
    log_error "ECK operator failed to become ready"
    exit 1
}

log_info "ECK operator is ready!"
sleep 10

# Step 5: Create Elasticsearch
log_info "Creating Elasticsearch cluster..."
kubectl apply -f "${K8S_DIR}/elasticsearch.yaml"

log_info "Waiting for Elasticsearch cluster to be ready (this may take a few minutes)..."
kubectl wait --for=condition=Ready pod -l common.k8s.elastic.co/type=elasticsearch -n elastic --timeout=600s || {
    log_warn "Elasticsearch pods not ready yet, waiting longer..."
}
sleep 15

# Step 6: Create Kibana
log_info "Creating Kibana..."
kubectl apply -f "${K8S_DIR}/kibana.yaml"

log_info "Waiting for Kibana to be ready (this may take a few minutes)..."
kubectl wait --for=condition=Ready pod -l common.k8s.elastic.co/type=kibana -n elastic --timeout=600s || {
    log_warn "Kibana pod not ready yet, continuing anyway..."
}
sleep 10

# Step 7: Create Service Account and RBAC for Fleet Server
log_info "Creating Fleet Server service account and RBAC..."
kubectl apply -f "${K8S_DIR}/fleet-server-rbac.yaml"

# Step 8: Create Fleet Server Agent
log_info "Creating Fleet Server agent..."
kubectl apply -f "${K8S_DIR}/fleet-server.yaml"

log_info "Waiting for Fleet Server to be ready (this may take a few minutes)..."
for i in {1..60}; do
    if kubectl get pod -l agent.k8s.elastic.co/name=fleet-server -n elastic 2>/dev/null | grep -q fleet-server; then
        kubectl wait --for=condition=Ready pod -l agent.k8s.elastic.co/name=fleet-server -n elastic --timeout=300s 2>/dev/null && break
    fi
    sleep 5
done
log_warn "Fleet Server pod check completed, continuing..."
sleep 10

# Step 9: Create Service Account and RBAC for Elastic Agent
log_info "Creating Elastic Agent service account and RBAC..."
kubectl apply -f "${K8S_DIR}/elastic-agent-rbac.yaml"

# Step 10: Create Elastic Agent DaemonSet
log_info "Creating Elastic Agent (DaemonSet)..."
kubectl apply -f "${K8S_DIR}/elastic-agent.yaml"

log_info "Waiting for Elastic Agent pods to be ready..."
for i in {1..60}; do
    if kubectl get pod -l agent.k8s.elastic.co/name=elastic-agent -n elastic 2>/dev/null | grep -q elastic-agent; then
        kubectl wait --for=condition=Ready pod -l agent.k8s.elastic.co/name=elastic-agent -n elastic --timeout=300s 2>/dev/null && break
    fi
    sleep 5
done
log_warn "Elastic Agent pods check completed, continuing..."

# Step 11: Create APM Service
log_info "Creating APM service..."
kubectl apply -f "${K8S_DIR}/apm-service.yaml"

# Step 12: Summary
log_info "========================================="
log_info "Deployment completed!"
log_info "========================================="

log_info "Checking deployment status..."
echo ""
echo "Namespaces:"
kubectl get namespaces | grep -E "elastic|apps"
echo ""
echo "Elasticsearch Status:"
kubectl get elasticsearch -n elastic
echo ""
echo "Kibana Status:"
kubectl get kibana -n elastic
echo ""
echo "Fleet Server Status:"
kubectl get agents -n elastic | grep fleet-server
echo ""
echo "Elastic Agent Status:"
kubectl get agents -n elastic | grep elastic-agent
echo ""
echo "Services:"
kubectl get svc -n elastic
echo ""
log_info "To access Kibana, use port-forward:"
log_info "kubectl port-forward -n elastic svc/kibana-sample-kb-http 5601:5601"
log_info "Then visit: http://localhost:5601"
echo ""
log_info "To get Elasticsearch and Kibana credentials:"
log_info "kubectl get secret -n elastic elasticsearch-sample-es-elastic-user -o jsonpath='{.data.elastic}' | base64 -d"
