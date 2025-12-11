# POC using elastic stack for logging, metrics, traces, APM


```shell
# fisrt create the cluster
terraform apply -target=kind_cluster.default -auto-approve
# second create eck
terraform apply -target=helm_release.eck -auto-approve

# then apply the rest of resources
terraform apply

# expose kibana
kubectl port-forward svc/kibana-sample-kb-http 5601:5601 -n elastic

# get kibana password for username elastic
kubectl get secret elasticsearch-sample-es-elastic-user -n elastic -o jsonpath='{.data.elastic}' | base64 -d && echo

http://localhost:5601

# build deploy the app
cd java-app
docker build -t java-app:latest .
kind load docker-image java-app:latest --name elastic-kind-cluster
kubectl apply -f deployment.yaml


kubectl run curl-test -n apps --image=curlimages/curl --rm -it --restart=Never -- sh -c 'for i in $(seq 1 10); do curl -s http://java-app:8080/api/hello; done'

```