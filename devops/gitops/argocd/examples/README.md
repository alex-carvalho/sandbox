

```shell
# create a k8s cluster
alias k=kubectl
kind  create cluster --name poc-argo

# create namespace to add the tests
k create namespace apps

# install argocd
k create namespace argocd
k apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# install argocd cli
VERSION=$(curl --silent "https://api.github.com/repos/argoproj/argo-cd/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')
sudo curl -sSL -o /usr/local/bin/argocd https://github.com/argoproj/argo-cd/releases/download/$VERSION/argocd-linux-amd64
sudo chmod +x /usr/local/bin/argocd

# foward the server api
k port-forward svc/argocd-server -n argocd 8080:443

# retrive initial admin argocd password
argocd admin initial-password -n argocd

# login on cli ignore certificate autosigned
argocd login localhost:8080 --insecure

# check list apps running, should be empty 
argocd app list


# install apps, we can use the argocd cli or apply a kube resource
# argocd app create guestbook --repo https://github.com/argoproj/argocd-example-apps.git --path guestbook --dest-server https://kubernetes.default.svc --dest-namespace default
k apply -f ./guestbook/argo-application.yaml

# now can list aplication on argcli 
argocd app list

# also we can list using kubectl
k get all -n apps
 
```

**UI**
https://localhost:8080/

