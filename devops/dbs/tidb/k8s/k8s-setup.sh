#!/usr/bin/env bash
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

kind create cluster --name "tidb-poc"

kubectl apply --server-side -f "${SCRIPT_DIR}/1-crd.yaml"

helm repo add pingcap https://charts.pingcap.com/ --force-update
kubectl get namespace tidb-admin >/dev/null 2>&1 || kubectl create namespace tidb-admin
helm upgrade --install --namespace tidb-admin tidb-operator pingcap/tidb-operator --version "v1.6.5" --wait --timeout 10m

kubectl wait --for=condition=Established --timeout=180s crd/tidbclusters.pingcap.com
kubectl wait --for=condition=Established --timeout=180s crd/tidbmonitors.pingcap.com
kubectl wait --for=condition=Established --timeout=180s crd/tidbdashboards.pingcap.com

# cluster
kubectl create namespace tidb-cluster
kubectl apply -n tidb-cluster -f "${SCRIPT_DIR}/2-cluster.yaml"
kubectl apply -n tidb-cluster -f "${SCRIPT_DIR}/3-monitor.yaml"
kubectl apply -n tidb-cluster -f "${SCRIPT_DIR}/4-dashboard.yaml"
kubectl apply -n tidb-cluster -f "${SCRIPT_DIR}/5-ngmonitor.yaml"


sleep 60

kubectl wait --for=condition=Ready pod -n tidb-admin --all --timeout=600s
kubectl wait --for=condition=Ready pod -n tidb-cluster --all --timeout=900s


echo "TiDB cluster and related components have been created successfully. You can access Grafana and TiDB Dashboard using port forwarding:"
nohup kubectl port-forward -n tidb-cluster svc/poc-grafana 3000 >/tmp/poc-grafana-port-forward.log 2>&1 &
nohup kubectl port-forward -n tidb-cluster svc/poc-tidb-dashboard-exposed 12333 >/tmp/poc-dashboard-port-forward.log 2>&1 &
nohup kubectl port-forward -n tidb-cluster svc/poc-tidb 4000 >/tmp/poc-tidb-port-forward.log 2>&1 &

open http://localhost:3000
open http://localhost:12333