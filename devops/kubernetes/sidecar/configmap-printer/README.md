# ConfigMap Printer

A simple Python application that reads a Kubernetes ConfigMap every 30 seconds.

## Environment Variables

- `POD_NAMESPACE` - Kubernetes namespace (default: `default`)
- `CONFIGMAP_NAME` - ConfigMap name to read (default: `properties`)

## Building add push the image to kind

```bash
./build-and-load.sh
```