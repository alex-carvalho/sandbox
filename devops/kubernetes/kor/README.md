# [kor](https://github.com/yonahd/kor)

A Golang Tool to discover unused Kubernetes Resources


kor can be installed on cluster as a cron job or deployment reporting to prometheus as well local

```shell
# install kor
kubectl krew install kor

# Apply resource
kubectl apply -f unused-resources.yaml
kubectl apply -f used-resources.yaml

# run kor all resources
kubectl kor all --show-reason --group-by=resource --output=table
```

kor output
```
kor version: v0.6.5

  _  _____  ____
 | |/ / _ \|  _ \
 | ' / | | | |_) |
 | . \ |_| |  _ <
 |_|\_\___/|_| \_\

Unused configmaps:
+---+-----------+---------------+-----------------------------------------------+
| # | NAMESPACE | RESOURCE NAME |                    REASON                     |
+---+-----------+---------------+-----------------------------------------------+
| 1 | default   | unused-config | ConfigMap is not used in any pod or container |
+---+-----------+---------------+-----------------------------------------------+

Unused services:
+---+-----------+----------------+-------------------------------+
| # | NAMESPACE | RESOURCE NAME  |            REASON             |
+---+-----------+----------------+-------------------------------+
| 1 | default   | unused-service | Service has no endpointslices |
+---+-----------+----------------+-------------------------------+

Unused secrets:
+---+-----------+---------------+------------------------------------------------------+
| # | NAMESPACE | RESOURCE NAME |                        REASON                        |
+---+-----------+---------------+------------------------------------------------------+
| 1 | default   | unused-secret | Secret is not used in any pod, container, or ingress |
+---+-----------+---------------+------------------------------------------------------+

Unused serviceaccounts:
+---+-----------+-----------------------+------------------------------+
| # | NAMESPACE |     RESOURCE NAME     |            REASON            |
+---+-----------+-----------------------+------------------------------+
| 1 | default   | unused-serviceaccount | ServiceAccount is not in use |
+---+-----------+-----------------------+------------------------------+

Unused :
+---+-----------+---------------+-------------------+
| # | NAMESPACE | RESOURCE NAME |      REASON       |
+---+-----------+---------------+-------------------+
| 1 | default   | unused-pvc    | PVC is not in use |
+---+-----------+---------------+-------------------+

Unused replicasets:
+---+-----------+---------------------+--------------------------+
| # | NAMESPACE |    RESOURCE NAME    |          REASON          |
+---+-----------+---------------------+--------------------------+
| 1 | default   | test-app-8687775766 | ReplicaSet is not in use |
+---+-----------+---------------------+--------------------------+
```