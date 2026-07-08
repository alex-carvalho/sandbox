#!/bin/bash
set -e

CLUSTER_NAME="${1:-flask-app}"
NAMESPACE="apps"
IMAGE_NAME="flask-app"
IMAGE_TAG="latest"

echo "🚀 Flask App Kind Cluster Setup"
echo "================================"

echo "🔍 Checking prerequisites..."
for cmd in kind docker helm kubectl; do
    if ! command -v $cmd &> /dev/null; then
        echo "❌ $cmd is not installed"
        exit 1
    fi
done

echo "✅ All prerequisites found"

echo ""
echo "1️⃣  Creating Kind cluster: $CLUSTER_NAME..."
if kind get clusters | grep -q "^$CLUSTER_NAME$"; then
    echo "⚠️  Cluster $CLUSTER_NAME already exists, skipping creation"
else
    kind create cluster --name "$CLUSTER_NAME"
    echo "✅ Kind cluster created"
fi

echo ""
echo "2️⃣  Setting kubeconfig context..."
kubectl cluster-info --context "kind-$CLUSTER_NAME"

echo ""
echo "3️⃣  Building Docker image..."
docker build -t "$IMAGE_NAME:$IMAGE_TAG" .
echo "✅ Docker image built"

echo ""
echo "4️⃣  Loading image into Kind cluster..."
kind load docker-image "$IMAGE_NAME:$IMAGE_TAG" --name "$CLUSTER_NAME"
echo "✅ Image loaded into Kind cluster"

echo ""
echo "5️⃣  Creating namespace..."
kubectl create namespace "$NAMESPACE" --dry-run=client -o yaml | kubectl apply -f -
echo "✅ Namespace ready"

echo ""
echo "6️⃣  Deploying with Helm..."
helm upgrade --install "$CLUSTER_NAME" ./helm \
    --namespace "$NAMESPACE" \
    --set image.repository="$IMAGE_NAME" \
    --set image.tag="$IMAGE_TAG" \
    --set image.pullPolicy="Never"

echo ""
echo "✅ Deployment complete!"

echo ""
echo "📊 Checking deployment status..."
kubectl get deployments -n "$NAMESPACE"
kubectl get pods -n "$NAMESPACE"

echo ""
echo "🔗 To access the app:"
echo "   kubectl port-forward -n $NAMESPACE svc/$CLUSTER_NAME 5000:80"
echo "   curl http://localhost:5000/health"

echo ""
echo "🧹 To cleanup:"
echo "   kind delete cluster --name $CLUSTER_NAME"
