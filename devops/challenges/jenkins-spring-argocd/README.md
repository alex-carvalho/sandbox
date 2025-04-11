## Create a Jenkins DSL that can deploy a simple Spring Boot app (build, run tests, run sonnar and deploy in K8S) use ArgoCD



```shell

alias k=kubectl

# run jenkins
wget https://get.jenkins.io/war-stable/2.492.2/jenkins.war
java -jar jenkins.war


# install kind
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.27.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# create cluster
kind create cluster --name kind-cluster


# install argocd
kubectl create namespace argocd
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# test pod service running on cluster
kubectl debug <pod> -it --image=nicolaka/netshoot
curl http://spring-web-3-j21/hello
```
