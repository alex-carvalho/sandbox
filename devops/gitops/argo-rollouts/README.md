# Argo Rollouts

**Argo Rollouts** is a Kubernetes controller that provides advanced deployment strategies such as blue-green, canary, and progressive delivery for Kubernetes applications. It extends the standard Kubernetes Deployment resource, enabling safer and more controlled application updates.

## What Problem Does It Solve?

Kubernetes Deployments offer basic rolling updates, but lack advanced deployment strategies needed for minimizing risk and enabling progressive delivery. Argo Rollouts solves this by providing:

- Fine-grained control over traffic shifting
- Automated analysis and promotion
- Easy rollback and pause/resume capabilities
- Integration with ingress controllers and service meshes

## Main Features

- **Canary Deployments:** Gradually shift traffic to new versions, with automated or manual promotion.
- **Blue-Green Deployments:** Deploy new versions alongside old ones, switch traffic instantly after verification.
- **Progressive Delivery:** Automated analysis and promotion based on metrics.
- **Pause/Resume:** Pause rollouts for manual checks or automated analysis.
- **Automated Rollbacks:** Roll back on failure or bad metrics.
- **Integration:** Works with Ingress, Service Meshes (Istio, Linkerd), and metric providers (Prometheus, Datadog, etc).

## Example: Canary Rollout

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: my-app
spec:
  replicas: 3
  strategy:
    canary:
      steps:
        - setWeight: 20
        - pause: {duration: 10m}
        - setWeight: 50
        - pause: {duration: 10m}
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
        - name: my-app
          image: my-app:v2
```

## Example: Blue-Green Rollout

```yaml
apiVersion: argoproj.io/v1alpha1
kind: Rollout
metadata:
  name: my-app
spec:
  replicas: 3
  strategy:
    blueGreen:
      activeService: my-app-active
      previewService: my-app-preview
      autoPromotionEnabled: false
  selector:
    matchLabels:
      app: my-app
  template:
    metadata:
      labels:
        app: my-app
    spec:
      containers:
        - name: my-app
          image: my-app:v2
```

## Instalation
1. Install Argo Rollouts:

   ```bash
   kubectl apply -f https://github.com/argoproj/argo-rollouts/releases/latest/download/install.yaml
   ```

2. Verify the rollout controller is running:

   ```bash
   kubectl get pods -n argo-rollouts
   ```

## References

- [Argo Rollouts Documentation](https://argoproj.github.io/argo-rollouts/)
- [GitHub Repository](https://github.com/argoproj/argo-rollouts)