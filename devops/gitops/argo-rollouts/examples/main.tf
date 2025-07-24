resource "kind_cluster" "default" {
  name = "poc-argo-rollouts"

  kind_config {
    raw = <<-EOT
      kind: Cluster
      apiVersion: kind.x-k8s.io/v1alpha4
      nodes:
        - role: control-plane
          extraPortMappings:
            - containerPort: 80
              hostPort: 80
              protocol: TCP
            - containerPort: 443
              hostPort: 443
              protocol: TCP
    EOT
  }

  wait_for_ready = true
}

resource "kubernetes_namespace" "argo_rollouts" {
  metadata {
    name = "argo-rollouts"
  }
}

resource "helm_release" "argo_rollouts" {
  name       = "argo-rollouts"
  namespace  = kubernetes_namespace.argo_rollouts.metadata[0].name
  repository = "https://argoproj.github.io/argo-helm"
  chart      = "argo-rollouts"
  version    = "2.40.2"
  depends_on = [kind_cluster.default]
}

resource "helm_release" "nginx_ingress" {
  name       = "ingress-nginx"
  namespace  = "ingress-nginx"
  repository = "https://kubernetes.github.io/ingress-nginx"
  chart      = "ingress-nginx"
  version    = "4.10.1"
  create_namespace = true
  depends_on = [kind_cluster.default]
  set {
    name  = "controller.nodeSelector.ingress-ready"
    value = "true"
  }
  set {
    name  = "controller.service.type"
    value = "NodePort"
  }
  set {
    name  = "controller.service.nodePorts.http"
    value = "80"
  }
  set {
    name  = "controller.service.nodePorts.https"
    value = "443"
  }
}