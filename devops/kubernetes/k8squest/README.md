
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
Fix the liveness probe configuration so pods stay running    

```shell
k edit  deploy api-server 
# update the liveness url removing the path
```

11 -  Pods receive traffic before they're ready, causing 502 errors for                                                                     ║
Add a readiness probe to prevent traffic from hitting pods too early 

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
Install metrics-server so HPA can read CPU/memory metrics and scale  
```shell
k apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml
# enable insecure tls flag
k patch deployment metrics-server -n kube-system --type='json' -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--kubelet-insecure-tls"}]'
```
