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

# install argocd cli
curl -sSL -o argocd-linux-amd64 https://github.com/argoproj/argo-cd/releases/latest/download/argocd-linux-amd64
sudo install -m 555 argocd-linux-amd64 /usr/local/bin/argocd
rm argocd-linux-amd64

# access argocd ui
kubectl port-forward svc/argocd-server -n argocd 8080:443

# get argocd initial password
argocd admin initial-password -n argocd


argocd app create spring-web-3-j21 --repo https://github.com/alex-carvalho/sandbox.git --path devops/challenges/jenkins-spring-argocd/k8s-artifacts --dest-server https://kubernetes.default.svc --dest-namespace default

# test pod service running on cluster
kubectl debug <pod> -it --image=nicolaka/netshoot
curl http://spring-web-3-j21/hello
```
