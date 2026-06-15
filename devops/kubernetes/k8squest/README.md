
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

22 - Ingress created but traffic doesn't reach backend service
- Fix the Ingress path rules to route traffic correctly 
```shell
k get ingress web-ingress -o yaml
k patch ingress web-ingress --type='json' -p='[{"op": "replace", "path": "/spec/rules/0/http/paths/0/path", "value": "/"}]'  
```

23 - A NetworkPolicy is blocking legitimate traffic between pods. The frontend can't reach the backend API, even though both pods are running and healthy. You need to identify and fix the overly restrictive network policy
- Fix the overly restrictive NetworkPolicy to allow frontend-to-backend communication
```shell
k get netpol backend-network-policy  -o yaml
k patch netpol backend-network-policy --type='json' -p='[{"op": "replace", "path": "/spec/ingress/0/from/0/podSelector/matchLabels/app", "value": "frontend"}]'
```

24 - A stateful application is randomly losing user sessions. Users log in successfully but their subsequent requests show them as logged out. The issue is that the Service is load-balancing requests across multiple pods, and session data isn't shared between them. You need to enable session affinity to fix this.
- Configure sessionAffinity on the Service to maintain user sessions
```shell
k patch service session-service --type='json' -p='[{"op": "replace", "path": "/spec/sessionAffinity", "value": "ClientIP"}]'
```

25 - Mission: A frontend application in the 'k8squest' namespace needs to call an API in the 'backend-ns' namespace, but it can't connect. The issue is using a short service name instead of the fully-qualified domain name (FQDN). You need to fix the DNS name to enable cross-namespace communication.
- Fix cross-namespace service communication using proper DNS FQDN
```shell
k get pod frontend-app -o yaml \
| sed 's/api-service/api-service.backend-ns/g' \
| kubectl replace --force -f -
```

26 - A service's endpoints aren't updating when pods restart or become unhealthy. Traffic continues routing to broken pods, causing errors. The issue is missing readiness probes—Kubernetes doesn't know when pods are actually ready to serve traffic. You need to add readiness probes to fix this.
- Add readiness probes so Service endpoints update correctly
```shell
k edit pod web-app-1
kubectl replace --force -f ....teporary file
```

29 - Your team deployed a service with type LoadBalancer, but it's been stuck in "Pending" state for 10 minutes. External clients can't access the application. The service works fine in the cloud environment, but your local development cluster can't provision the LoadBalancer. You need to understand the difference between service types and choose the right one for local development.                      ║
- Fix the broken resources and make the validation pass
```shell
kubectl patch svc web-service \
  --type=merge \
  -p '{
    "spec":{
      "type":"NodePort",
      "ports":[
        {
          "port":80,
          "targetPort":80
        }
      ]
    }
  }'
```

30 - Your StatefulSet pods can't communicate with each other using predictable DNS names. The application expects to reach individual pods at pod-0.service-name, pod-1.service-name, but DNS resolution isn't working. The service is configured with a ClusterIP, but StatefulSets require a special type of service called a "headless service" for direct pod-to-pod DNS resolution.
- Fix the broken resources and make the validation pass, make web-cluster a headless service, ClusterIP: None
```shell
k get svc web-cluster -o yaml > /tmp/svc.yaml 
vi /tmp/svc.yaml 
k replace --force -f /tmp/svc.yaml
```

31 - Your application pod is stuck in ContainerCreating state because its PersistentVolumeClaim is in Pending status. The pod needs persistent storage to save database files, but the storage isn't being provisioned. You need to investigate why the PVC can't bind to a PersistentVolume and fix the storage configuration.
- Fix the broken resources and make the validation pass
```shell
k get pvc app-storage-claim -o yaml
kubectl patch pv app-storage \
  -p '{"spec":{"capacity":{"storage":"5Gi"},"storageClassName":"fast"}}'
```

32 - Your application pod is crashing with errors about missing files. The pod expects to read configuration from /app/config, but the volume is mounted at the wrong path. You need to fix the volumeMount configuration to mount the storage at the correct location so the application can find its config files.
Fix the broken resources and make the validation pass
```shell  
k get pod web-app  -o yaml > pod.yaml
vi pod.yaml # remove status part ( :/^status:/,$d) and annotations, fix volumeMount path to /app/config
k replace --force -f pod.yaml
```

33 - Multiple pods need to share the same storage volume across nodes, but the storage is configured with ReadWriteOnce instead of ReadWriteMany. This level runs on Kind (single-node), so all 3 pods will appear to work. However, the configuration is WRONG for production multi-node clusters. The validation will detect this, and you'll also learn that PVC specs are immutable when you try to fix it
This teaches TWO lessons: (1) correct access mode selection, (2) PVC immutability.
Fix the broken resources and make the validation pass
```shell

k get pvc shared-pvc -o yaml > /tmp/pvc.yaml
vi /tmp/pvc.yaml # change access mode to ReadWriteMany
k replace --force -f /tmp/pvc.yaml
k get pv shared-storage -o yaml > /tmp/pv.yaml
vi /tmp/pv.yaml # change access mode to ReadWriteMany
k replace --force -f /tmp/pv.yaml
```

34 - A StatefulSet for a database cluster has all pods writing to the same PVC! Each pod should have its own dedicated storage, but they're all sharing one volume. This is causing data corruption and conflicts. 
- Fix the broken resources and make the validation pass
```shell
k get pvc
k get sts postgres-cluster -o yaml > sts.yaml
vi sts.yaml
# add to spec of statefulset and remove pod level volumeClaimRef
volumeClaimTemplates:
- metadata:
    name: database-storage
  spec:
    accessModes: [ReadWriteOnce]
    resources:
      requests:
        storage: 5Gi
    storageClassName: standard

k replace --force -f sts.yaml
```

35 - A PVC has been sitting in Pending state for hours. The pod can't start because it's waiting for storage. The PVC references a StorageClass that doesn't exist!
- Fix the broken resources and make the validation pass
```shell
k get pvc -o wide
k get pv app-storage -o yaml > /tmp/pv.yaml
vi /tmp/pv.yaml # change storageClassName to standard
k replace --force -f /tmp/pv.yaml
```
