# 02-hpa-prometheus: Autoscaling with Standard HPA and Prometheus Custom Metrics

This module demonstrates autoscaling Kubernetes workloads using **standard Kubernetes HPA (autoscaling/v2)** with custom metrics exposed through **Prometheus** and the **Prometheus Adapter**, connected to a Redis list named `hpa-prometheus-queue`.

## Key Concepts

- **Metric Ingestion Pipeline**:
  `Redis (LLEN)` $\to$ `redis_exporter (:9121)` $\to$ `Prometheus scrape (every 5s)` $\to$ `Prometheus Adapter (custom.metrics.k8s.io)` $\to$ `Kubernetes HPA controller`.
- **Native Scale-to-Zero (Kubernetes v1.37+)**:
  Starting in Kubernetes v1.37, `minReplicas: 0` is officially supported for `External` and `Object` metrics. Both KEDA and native HPA can now scale workloads down to 0 replicas.
- **Latency Delay (Cold-Start Contrast)**:
  Even though native HPA can scale to 0 in v1.37+, it still relies on the multi-hop metric pipeline (scrape intervals + Adapter polling). Waking up from 0 replicas takes **~20–35s** with HPA + Prometheus, compared to **~5–10s** with KEDA's direct Redis polling.

## Manifests

- `manifests/namespace.yaml`: Creates the `hpa-poc` namespace.
- `manifests/consumer.yaml`: Consumer deployment pulling from `hpa-prometheus-queue` via `BLPOP`. Baseline `replicas: 0`.
- `manifests/hpa.yaml`: Standard `autoscaling/v2` HPA resource:
  - Target: `hpa-consumer`
  - `minReplicas: 0` (v1.37+ native scale-to-zero)
  - `maxReplicas: 8`
  - Metric: `External` metric `hpa_prometheus_queue_length` with `averageValue: "5"`.
  - Behavior: Fast scale-down stabilization window (20s) for testing.

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
1. Confirm baseline pod is at **1** (cannot scale to zero).
2. Inject 30 messages into `hpa-prometheus-queue`.
3. Display latency while the metric traverses the Prometheus exporter and adapter pipeline.
4. Display scale out to **6 pods**.
5. Watch workers drain the queue.
6. Watch HPA scale back down and stop at **1 pod**.

### 3. Clean up
```bash
./cleanup.sh
```
