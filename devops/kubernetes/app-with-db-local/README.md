# Kubernetes

__Configuration for Ingress:__

https://kubernetes.github.io/ingress-nginx/deploy/#prerequisite-generic-deployment-command

```
kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/mandatory.yaml

minikube addons enable ingress

kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/master/deploy/provider/baremetal/service-nodeport.yaml
```

__Add db password as secret:__
```
kubectl create secret generic todolistdbpass --from-literal PGPASSWORD=secret
``` 


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
