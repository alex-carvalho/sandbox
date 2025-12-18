# Go Sidecar - Properties File ConfigMap Manager

A Kubernetes sidecar container written in Go that syncs properties from a file into a ConfigMap.

## Overview

This sidecar reads a properties file and updates a Kubernetes ConfigMap with its contents. It's designed to be:
- **One-shot execution**: Reads once, updates ConfigMap, exits cleanly
- **Kubernetes-aware**: Auto-detects its namespace from service account
- **Simple and lightweight**: Single file input, no file watching or polling

## Configuration

The sidecar is configured via environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `PROPERTIES_FILE` | `/etc/config/application.properties` | Path to properties file to read |
| `POD_NAMESPACE` | Auto-detect* | Kubernetes namespace for ConfigMap |
| `CONFIGMAP_NAME` | `properties` | ConfigMap resource name |

*Auto-detects from `/var/run/secrets/kubernetes.io/serviceaccount/namespace` when running in Kubernetes


## Installation & Deployment

### 1. Build Docker Image

```bash
make docker-build
```

## Usage

### Quick Test

Test locally with your Kubernetes cluster:

```bash
# Creates sample properties file and updates ConfigMap
make local-test

# With custom file
PROPERTIES_FILE=/path/to/file.properties make local-test
```

