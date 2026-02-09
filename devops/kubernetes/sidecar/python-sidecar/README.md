# Python Sidecar - Properties File ConfigMap Manager

A Kubernetes sidecar container written in Python that watches a properties file and syncs changes to a ConfigMap.

## Overview

This sidecar continuously watches a properties file directory for changes and automatically syncs updates to a Kubernetes ConfigMap. It's designed to be:
- **File watcher**: Detects file changes in real-time using watchdog
- **ConfigMap syncer**: Automatically updates ConfigMap when properties change
- **Kubernetes-aware**: Auto-detects its namespace from service account
- **Lightweight**: Efficient file watching with minimal resource usage

## Configuration

The sidecar is configured via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `WATCH_DIR` | `/etc/config` | Directory to watch for changes |
| `POD_NAMESPACE` | `default` | Kubernetes namespace for ConfigMap |
| `CONFIGMAP_NAME` | `properties` | ConfigMap resource name |

## Usage

```bash

# Deploy to Kubernetes (builds, loads to kind, and deploys)
make deploy

# View logs
kubectl logs -f deployment/app-with-python-sidecar -c python-sidecar
kubectl logs -f deployment/app-with-python-sidecar -c configmap-printer


# Add a property to test file watching
POD=$(kubectl get pods -l app=python-app -o jsonpath='{.items[0].metadata.name}')
kubectl exec -it $POD -c python-sidecar -- sh -c 'echo "test.property=value" >> /etc/config/application.properties'

```