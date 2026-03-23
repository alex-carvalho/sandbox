
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