# Ansible + Kubernetes

use of kubernetes.core


## Collections Used

| Collection | Modules used |
|---|---|
| `kubernetes.core` | `k8s` — create/update any k8s resource |
| | `k8s_info` — query resources and wait for readiness |

## Usage

```bash
make cluster
make init   
make run    
make check  
```

## Testing

```bash
curl http://localhost:8080
```
