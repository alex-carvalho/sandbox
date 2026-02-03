#!/bin/bash

set -e

echo "Building Docker image..."
docker build -t configmap-printer:latest .

echo "Loading image into kind cluster..."
kind load docker-image configmap-printer:latest

echo "✓ Image built and loaded to kind cluster"
echo ""
echo "Deploy with: kubectl apply -f deployment.yaml"
