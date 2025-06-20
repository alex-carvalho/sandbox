# ArgoCD

- GitOps delivery tool
- Git souce of true
- Argo deploy the app to k8s cluster
- Keep the desired state in git repository in sync with kubernetes
- Pull based

## ArgoCD Architecture
![ArgoCD Architecture](argocd-architecture.png)

### Components

__API Server__
- gRPC/REST server which exposes the API consumed by the Web UI, CLI, and CI/CD systems

__Repository Server__
- local cache of the Git repository

__Application Controller__
- is a Kubernetes controller which continuously monitors running applications and compares the current


### ArgoCD Architecture Component
![ArgoCD Architecture Component](argocd-architecture-component.png)