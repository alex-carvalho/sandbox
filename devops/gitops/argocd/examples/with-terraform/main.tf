resource "kind_cluster" "default" {
    name = "poc-argocd-cluster"
}

resource "kubernetes_namespace" "apps" {
    metadata {
        name = "apps"
    }
}

resource "kubernetes_namespace" "argocd" {
    metadata {
        name = "argocd"
    } 
}