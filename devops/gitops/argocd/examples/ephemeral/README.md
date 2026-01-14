# using argocd no git source

### cleanup the cluster
kubectl delete namespace ephemeral-apps
kubectl delete namespace argocd
kubectl delete namespace apps
kubectl delete clusterrole argocd-application-controller argocd-notifications-controller argocd-server --ignore-not-found
kubectl delete clusterrolebinding argocd-application-controller argocd-notifications-controller argocd-server --ignore-not-found


# add helm repo:
helm repo add argo https://argoproj.github.io/argo-helm
helm repo update


# Install ArgoCD with the values file:
helm install argocd argo/argo-cd \
  --namespace argocd \
  --create-namespace \
  -f argocd-values.yaml

# Create the ephemeral-apps namespace
k create namespace ephemeral-apps

# Create RBAC rolebinding for ArgoCD controller
kubectl create rolebinding argocd-ephemeral \
  -n ephemeral-apps \
  --clusterrole=admin \
  --serviceaccount=argocd:argocd-application-controller

# Apply the AppProject and Applications
k apply -f ./project.yml
k apply -f ./grafana.yml


# checking
### retrive initial admin argocd password
argocd admin initial-password -n argocd

### foward the server api
k port-forward svc/argocd-server -n argocd 8080:443

