
kind create cluster

kubectl create -f https://raw.githubusercontent.com/pingcap/tidb-operator/v1.6.5/manifests/crd.yaml


helm repo add pingcap https://charts.pingcap.com/
kubectl create namespace tidb-admin
helm install --namespace tidb-admin tidb-operator pingcap/tidb-operator --version v1.6.5

kubectl get pods --namespace tidb-admin


# cluster
kubectl create namespace tidb-cluster
kubectl -n tidb-cluster apply -f https://raw.githubusercontent.com/pingcap/tidb-operator/v1.6.5/examples/basic/tidb-cluster.yaml




# Deploy TiDB Dashboard independently
kubectl -n tidb-cluster apply -f https://raw.githubusercontent.com/pingcap/tidb-operator/v1.6.5/examples/basic/tidb-dashboard.yaml

#Deploy TiDB monitoring services
kubectl -n tidb-cluster apply -f https://raw.githubusercontent.com/pingcap/tidb-operator/v1.6.5/examples/basic/tidb-monitor.yaml


kubectl get po -n tidb-cluster


# grafana (admin/admin)
kubectl port-forward -n tidb-cluster svc/basic-grafana 3000


# tidb dashboard (root/"")
kubectl port-forward -n tidb-cluster svc/basic-tidb-dashboard-exposed 12333 