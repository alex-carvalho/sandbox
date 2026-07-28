# Loki & OpenTelemetry Collector Observability POC

This repository contains a Proof of Concept (POC) to deploy a local observability stack using a 3-node Kind Kubernetes cluster. It showcases how the **OpenTelemetry (OTel) Collector** collects container logs using the `filelog` receiver, enriches them with Kubernetes metadata using the `k8s_attributes` processor, and exports them to **Grafana Loki**.

## Prerequisites

Ensure you have the following CLI tools installed:
- [kind](https://kind.sigs.k8s.io/)
- [kubectl](https://kubernetes.io/docs/reference/kubectl/)
- [helm](https://helm.sh/)

## Running the POC

Simply run the setup script to provision the cluster and deploy Loki, Grafana, OpenTelemetry Collector, and a demo log generator application:

```bash
./setup.sh
```

## Accessing the Dashboard

Once all components are deployed and ready, port-forward to Grafana:

```bash
kubectl port-forward svc/grafana 3000:80 -n observability
```

- **URL**: [http://localhost:3000](http://localhost:3000)
- **Username**: `admin`
- **Password**: `admin`
