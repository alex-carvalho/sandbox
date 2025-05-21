
variable "prefix" {
  type    = string
  default = ""
  description = "Used on tests to change the name of the resources"
}

resource "kind_cluster" "default" {
    name = "${var.prefix}poc-kind-cluster"
}

resource "kubernetes_namespace" "apps" {
    metadata {
        name = "apps"
    }
}

resource "kubernetes_namespace" "infra" {
    metadata {
        name = "infra"
    }
}

resource "helm_release" "prometheus" {
  name       = "prometheus"
  namespace  = kubernetes_namespace.infra.metadata[0].name
  repository = "https://prometheus-community.github.io/helm-charts"
  chart      = "prometheus"
  version    = "27.14.0"

  create_namespace = false
}

resource "helm_release" "grafana" {
  name       = "grafana"
  namespace  = kubernetes_namespace.infra.metadata[0].name
  repository = "https://grafana.github.io/helm-charts"
  chart      = "grafana"
  version    = "9.0.0"

  create_namespace = false

  values = [
    file("${path.module}/grafana-values.yaml")
  ]

  set {
    name  = "service.type"
    value = "ClusterIP"
  }
}