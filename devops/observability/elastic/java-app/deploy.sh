#!/bin/bash

set -e

docker build -t java-app:latest .
kind load docker-image java-app:latest --name elastic-kind-cluster
kubectl apply -f deployment.yaml