

Build jenkins agent image: 
```shell
docker build -t ghcr.io/alex-carvalho/sandbox/jenkins-agent:latest -f jenkins-agent.Dockerfile .
```

Push the jenkins agent image:
```shell
docker push ghcr.io/alex-carvalho/sandbox/jenkins-agent:latest
```

kubectl port-forward svc/jenkins -n infra 8080:8080


Go to infra and run `tofu apply -auto-approve` it will create a k8s kind cluster with two namespaces (apps, infra) and install grafana and prometheus on it, to access Grafana run  

```shell
kubectl port-forward -n infra svc/grafana 3000:80
```
