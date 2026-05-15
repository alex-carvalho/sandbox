
# k8squest 
From https://github.com/Manoj-engineer/k8squest

Shortcut for easy quest:
``` shell
kubectl config set-context --current --namespace=k8squest
```

<br>
<br>

1 - Pod CrashLoopBackOff wrong command
- You can't edit most fields of a running pod, need to recreate it
``` shell
k get pod -o yaml  nginx-broken > my-pod.yaml
vi my-pod.yaml 
k delete pod nginx-broken         
k apply -f my-pod.yaml    
```

2 - Deployment has zero replicas
- Simple doing using imperative way, or can edit the yaml

```shell
k scale deployment.apps/web --replicas=1 # or  k edit deployment.apps/web 
```

3 - A pod is stuck trying to pull a container image that doesn't exist 
- Similar to num 1, need to recreate the pos to change the tag
``` shell
k get pod -o yaml web-app > my-pod.yaml
vi my-pod.yaml 
k delete pod web-app        
k apply -f my-pod.yaml    
```

4 - A pod is stuck in Pending status because it's asking for more resources than available   
- Similat to num 1, pod cannot rquest more resource that available by a single node in the cluster
``` shell
k get pod -o yaml hungry-app > my-pod.yaml
vi my-pod.yaml 
k delete pod hungry-app       
k apply -f my-pod.yaml    
```

4 - Service cannot reach pod, fix label 
- Can edit the service or patch it directly 
```shell
k patch svc backend-service -p '{"spec":{"selector":{"app":"backend"}}}'
```

5 - Service port does not match pod port
- Can edit the service or patch the targetPort to match the pod port

```shell
k patch svc web-service -p '{"spec":{"ports":[{"port":80,"targetPort":80}]}}'
```

6 - A multi-container pod where the sidecar container keeps crashing
- The file that log sidecar container was reading was not created yet, just use `-F` on tail command fix
```shell
k get pod -o yaml  app-with-logging   > my-pod.yaml
vi my-pod.yaml 
k delete pod app-with-logging           
k apply -f my-pod.yaml    
```

7 - Pod is running but the application isn't working. The problem is only visible in logs 
- Check logs pod, edit to add missing env
```shell
k logs database-app
k get pod -o yaml  database-app > my-pod.yaml
vi my-pod.yaml 
k delete pod database-app          
k apply -f my-pod.yaml   
```

8 -  Pod stuck in Init state because the init container never completes        
- Need to move the pod and the service to k8squest namespace
```shell
k get service -o yaml -n default backend-service  > my-service.yaml
k get pod -o yaml -n default client-app  > my-pod.yaml
vi my-service.yaml
vi my-pod.yaml 
k delete pod -n default client-app  
k delete service -n default backend-service     
k apply -f my-pod.yaml  
k apply -f my-service.yaml 
```

9 - A rolling update has failed because the new container image doesn't exist   
- Need rollback the deployment to the previous working version    

```shell
k get rs 
k rollout status deployment/web-app
k rollout undo deployment/web-app
```

10 -  Pods keep restarting because the liveness probe is checking the wrong endpoint   
- Fix the liveness probe configuration so pods stay running    

```shell
k edit  deploy api-server 
# update the liveness url removing the path
```

11 -  Pods receive traffic before they're ready, causing 502 errors for                                                                     ║
- Add a readiness probe to prevent traffic from hitting pods too early 

```shell
k edit deployment slow-startup-app
readinessProbe:
  httpGet:
    path: /
    port: 80
  initialDelaySeconds: 22  # add sleep greater than sleep on start
  periodSeconds: 5
  successThreshold: 1
```

12 - HorizontalPodAutoscaler fails to scale pods due to missing metrics
- Install metrics-server so HPA can read CPU/memory metrics and scale  
```shell
k apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
# enable insecure tls flag
k patch deployment metrics-server -n kube-system --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'
```

13 - Rolling update with maxUnavailable: 100% causes complete service 
- Fix the rollout strategy to prevent all pods from being down simultaneously     
```shell
k edit deploy critical-api
# update maxUnavailable to 0 and maxSurge to 1
```

 14 - PodDisruptionBudget is too restrictive, preventing node drains and updates 
 - Fix the PDB to allow safe pod evictions while maintaining availability   
 ```shell
 k get pdb db-proxy-pdb  -o yaml
 # update the minAvailable value to allow safe pod evictions
 k patch pdb db-proxy-pdb --type='json' -p='[{"op": "replace", "path": "/spec/minAvailable", "value": 1}]'
 ```

15 - Service is pointing to the old version, users don't see the new deployment
- Update the service selector to route traffic to the new blue-green deployment
```shell
k get service app-service -o yaml
k patch service app-service  -p '{"spec":{"selector":{"version":"green"}}}'
```

16 - Canary deployment has wrong replica ratios, making testing
- Fix the replica counts to achieve proper 90/10 traffic split for canary testing     
```shell
k scale deployment app-stable --replicas=9
k scale deployment app-canary --replicas=1
```

17 - Using Deployment for a database causes data loss and pod identity issues                                                                                                                                 ║
- Convert the Deployment to a StatefulSet to provide stable pod identities and persistent storage  
```shell
k delete deployment database
k apply -f - <<EOF
apiVersion: apps/v1
kind: StatefulSet
metadata:
  name: database
  namespace: k8squest
spec:
  serviceName: "database-service"  # Important!
  replicas: 3
  selector:
    matchLabels:
      app: database
  template:
    metadata:
      labels:
        app: database
    spec:
      containers:
      - name: db
        image: hashicorp/http-echo:latest
        args:
          - "-text=Database Pod"
          - "-listen=:8080"
        ports:
        - containerPort: 8080
EOF

```

18 - Manual ReplicaSet makes updates difficult and breaks automation
- Convert the standalone ReplicaSet to a Deployment for better management
```shell
kubectl get rs web-app-rs -o yaml \
| sed 's/kind: ReplicaSet/kind: Deployment/' \
| sed 's/name: web-app-rs/name: web-app/' \
| sed 's/image: hashicorp/http-echo:0.2.3/image: hashicorp/http-echo:latest/' \
| grep -v 'resourceVersion:\|uid:\|selfLink:\|creationTimestamp:\|generation:\|managedFields:\|ownerReferences:\|fullyLabeledReplicas:' \
> deployment.yaml
k delete replicaset web-app-rs
k apply -f deployment.yaml

```

19 - A service exists but pods aren't receiving traffic
- Fix the service so it routes traffic to the backend pods
```shell
k get service backend-service -o yaml
k patch service backend-service  -p '{"spec":{"selector":{"app":"backend"}}}'
```

20 - NodePort service created but can't access from outside cluster
- Fix the NodePort service configuration to enable external access
```shell
k get svc web-nodeport -o yaml   
k patch svc web-nodeport --type='json' -p='[{"op": "replace", "path": "/spec/ports/0/nodePort", "value": 30080}]'
```

21 - Pods can't resolve service names via DNS
- Fix DNS configuration so pods can resolve service names
```shell
k logs app-client
kubectl get pod app-client -o yaml \
| sed 's/-h database/-h database-service/g' \
| kubectl replace --force -f -
```
