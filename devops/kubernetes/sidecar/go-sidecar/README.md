# Go Sidecar - Properties File ConfigMap Manager

A Kubernetes sidecar container written in Go that watches a properties file and syncs changes to a ConfigMap.

## Overview

This sidecar continuously watches a properties file directory for changes and automatically syncs updates to a Kubernetes ConfigMap. It's designed to be:
- **File watcher**: Detects file changes in real-time using fsnotify
- **ConfigMap syncer**: Automatically updates ConfigMap when properties change
- **Kubernetes-aware**: Auto-detects its namespace from service account
- **Lightweight**: Efficient file watching with minimal resource usage

## Configuration

The sidecar is configured via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PROPERTIES_FILE` | `/etc/config/application.properties` | Path to properties file |
| `WATCH_DIR` | `/etc/config` | Directory to watch for changes |
| `POD_NAMESPACE` | Auto-detect* | Kubernetes namespace for ConfigMap |
| `CONFIGMAP_NAME` | `properties` | ConfigMap resource name |

*Auto-detects from `/var/run/secrets/kubernetes.io/serviceaccount/namespace` when running in Kubernetes


## Usage

```bash
# Creates sample properties file and updates ConfigMap locally
make local-test

# build image docker, load to kind cluster and deploy the app with sidecar
deploy deploy


# check logs
POD=$(kubectl get pods -l app=myapp -o jsonpath='{.items[0].metadata.name}')

kubectl logs deployment/app-with-sidecar -c go-sidecar 
kubectl logs deployment/app-with-sidecar -c configmap-printer

# modify the file and check the logs, the watcher will thrigger file changes
kubectl exec -it $POD -c go-sidecar -- sh -c 'echo "new.property=test-value" >> /etc/config/application.properties'
```

