
terraform {
  required_providers {
    argocd = {
      source  = "argoproj-labs/argocd"
      version = "7.8.2"
    }
    kind = {
      source  = "tehcyx/kind"
      version = "~> 0.2.1"
    }
    kubernetes = {
      source  = "hashicorp/kubernetes"
      version = "~> 2.37.1"
    }

    helm = {
      source  = "hashicorp/helm"
      version = "~> 3.0.2"
    }
  }
}

provider "kind" {}

provider "kubernetes" {
    config_path = kind_cluster.default.kubeconfig_path
}

provider "helm" {
    kubernetes = {
        config_path = kind_cluster.default.kubeconfig_path
    }
}

resource "helm_release" "argocd" {
    name       = "argocd"
    repository = "https://argoproj.github.io/argo-helm"
    chart      = "argo-cd"
    version    = "5.51.6"
    namespace  = "argocd"

    create_namespace = true

    values = [
    ]
}

data "kubernetes_service" "argocd_server" {
  metadata {
    name      = "argocd-server"
    namespace = "argocd"
  }
  depends_on = [helm_release.argocd]
}

data "kubernetes_secret" "argocd_admin" {
  metadata {
    name      = "argocd-initial-admin-secret"
    namespace = "argocd"
  }
  depends_on = [helm_release.argocd]
}


provider "argocd" {
  server_addr = "localhost:8080"
  username    = "admin"
  password    = data.kubernetes_secret.argocd_admin.data["password"]
  insecure    = true
}


# output "argocd_initial_admin_secret" {
#   value = nonsensitive(data.kubernetes_secret.argocd_admin.data["password"])
# }

 