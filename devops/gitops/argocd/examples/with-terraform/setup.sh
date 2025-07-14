kubectl wait deployment/argocd-server -n argocd --for=condition=available --timeout=120s
kubectl port-forward svc/argocd-server -n argocd 8080:443