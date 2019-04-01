# Kubernetes

## __Test app__

- deploy
```
kubectl apply -f . 
```

- check pods
```
kubectl get pods 
```

- test
```
curl -v http://$(minikube ip):31388/actuator/health
```

- delete
```
kubectl delete -f .
```
