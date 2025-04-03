## Create a Jenkins DSL that can deploy a simple Spring Boot app (build, run tests, run sonnar and deploy in k8s) use HELM

```shell

alias k=kubectl

# run jenkins
wget https://get.jenkins.io/war-stable/2.492.2/jenkins.war
java -jar jenkins.war

# run docker registry local
docker run -d --restart=always -p 5000:5000 --name my-registry registry:2


# install kind
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.27.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# create cluster
kind create cluster --name kind-cluster --config kind-config.yaml

# conect the registry with kind
docker network connect kind my-registry

# Install helm
curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# test pod service running on cluster
kubectl debug <pod> -it --image=nicolaka/netshoot
curl http://spring-web-3-j21:8080/hello
```
