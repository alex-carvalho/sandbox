terraform {
  required_providers {
    kind = {
      source  = "tehcyx/kind"
      version = "~> 0.2.1"
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
