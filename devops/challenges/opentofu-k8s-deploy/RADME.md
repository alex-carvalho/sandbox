

kubectl port-forward -n infra svc/grafana 3000:80


service:
    metadata:
        annotations:
            prometheus.io/scrape: "true"
            prometheus.io/path: /actuator/prometheus
            prometheus.io/port: "8080"
