

Go to infra and run `tofu apply -auto-approve` it will create a k8s kind cluster with two namespaces (apps, infra) and install grafana and prometheus on it, to access Grafana run  `kubectl port-forward -n infra svc/grafana 3000:80`
