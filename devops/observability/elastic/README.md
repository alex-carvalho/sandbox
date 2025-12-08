# POC using elastic stack for logging, metrics, traces, APM


```shell
# apply resources
terraform apply

# expose kibana
kubectl port-forward svc/kibana-sample-kb-http 5601:5601 -n elastic

# get kibana password for username elastic
kubectl get secret elasticsearch-sample-es-elastic-user -n elastic -o jsonpath='{.data.elastic}' | base64 -d && echo

http://localhost:5601

# build deploy the app
cd java-app
docker build -t java-app:latest .
kind load docker-image java-app:latest --name elastic-ki
kubectl apply -f deployment.yaml
```