# POC using elastic stack for logging, metrics, traces, APM


```shell
./deploy.sh

# get kibana password for username elastic
kubectl get secret elasticsearch-sample-es-elastic-user -n elastic -o jsonpath='{.data.elastic}' | base64 -d && echo

# expose kibana
kubectl port-forward svc/kibana-sample-kb-http 5601:5601 -n elastic

http://localhost:5601


kubectl run curl-test -n apps --image=curlimages/curl --rm -it --restart=Never -- sh -c 'for i in $(seq 1 10); do curl -s http://java-app:8080/api/hello; done'

```