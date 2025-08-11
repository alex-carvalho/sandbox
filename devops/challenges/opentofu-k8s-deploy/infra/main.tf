
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

resource "helm_release" "jenkins" {
  name             = "jenkins"
  namespace        = kubernetes_namespace.infra.metadata[0].name
  repository       = "https://charts.jenkins.io"
  chart            = "jenkins"
  version          = "5.8.73"
  create_namespace = false
  
  values = [
    file("${path.module}/jenkins-values.yaml")
  ]
  
  depends_on = [kind_cluster.default]
}