# 01-keda: Autoscaling with KEDA (Kubernetes Event-driven Autoscaling)

This module demonstrates autoscaling Kubernetes workloads using **KEDA** connected directly to a Redis list named `keda-queue`.

## Key Concepts

- **Direct Metric Polling**: KEDA queries the Redis broker directly (`LLEN keda-queue`) every 5 seconds. It does not require Prometheus, an exporter, or an external metrics pipeline.
- **Scale-to-Zero (Cold Start)**: When `keda-queue` is empty, KEDA scales the consumer deployment down to **0 replicas**. No CPU or RAM is consumed when idle.
- **Immediate Activation**: When messages arrive, KEDA immediately activates the deployment from 0 to 1 replica and dynamically manages the HPA to scale out to the required target count.

## Manifests

- `manifests/namespace.yaml`: Creates the `keda-poc` namespace.
- `manifests/consumer.yaml`: Consumer deployment pulling from `keda-queue` via `BLPOP`. Defaults to `replicas: 0`.
- `manifests/scaledobject.yaml`: The KEDA `ScaledObject` CRD configuring:
  - Target: `keda-consumer`
  - `minReplicaCount: 0`
  - `maxReplicaCount: 8`
  - `pollingInterval: 5` (seconds)
  - `cooldownPeriod: 20` (seconds)
  - Trigger: `redis` scaler with `listName: keda-queue` and `listLength: "5"`.

## Usage

### 1. Deploy
```bash
./deploy.sh
```

### 2. Run Test Demo
```bash
./test.sh
```

The test script will:
1. Confirm initial pods are at **0**.
2. Inject 30 messages into `keda-queue`.
3. Display real-time scaling out to **6 pods** (`30 msgs / 5 = 6 pods`).
4. Watch workers drain the queue.
5. Watch KEDA cooldown and scale back down to **0 pods**.

### 3. Clean up
```bash
./cleanup.sh
```
