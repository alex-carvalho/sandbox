
resource "kind_cluster" "default" {
  name = "elastic-kind-cluster"
}


resource "kubernetes_namespace" "apps" {
  metadata {
    name = "apps"
  }
}

resource "helm_release" "eck" {
  name             = "eck-operator"
  repository       = "https://helm.elastic.co"
  chart            = "eck-operator"
  version          = "2.16.0"
  namespace        = "elastic-system"
  create_namespace = true
  wait             = true
  timeout          = 600
  skip_crds        = false

  values = [
    <<-EOF
    operator:
      watchNamespaces: []
    EOF
  ]
}

resource "kubernetes_namespace" "elastic" {
  metadata {
    name = "elastic"
  }
}


resource "kubernetes_manifest" "es_cr" {
  manifest = {
    apiVersion = "elasticsearch.k8s.elastic.co/v1"
    kind       = "Elasticsearch"
    metadata = {
      name      = "elasticsearch-sample"
      namespace = kubernetes_namespace.elastic.metadata[0].name
    }
    spec = yamldecode(file("${path.module}/k8s/elasticsearch.yaml")).spec
  }
  
  depends_on = [
    kubernetes_namespace.elastic,
    helm_release.eck
  ]
}

# Kibana
resource "kubernetes_manifest" "kibana_cr" {
  manifest = {
    apiVersion = "kibana.k8s.elastic.co/v1"
    kind       = "Kibana"
    metadata = {
      name      = "kibana-sample"
      namespace = kubernetes_namespace.elastic.metadata[0].name
    }
    spec = yamldecode(file("${path.module}/k8s/kibana.yaml")).spec
  }
  
  depends_on = [kubernetes_manifest.es_cr]
}

# Service Account for Fleet Server
resource "kubernetes_service_account" "fleet_server" {
  metadata {
    name      = "fleet-server"
    namespace = kubernetes_namespace.elastic.metadata[0].name
  }
}

resource "kubernetes_cluster_role" "fleet_server" {
  metadata {
    name = "fleet-server"
  }

  rule {
    api_groups = [""]
    resources  = ["pods", "namespaces", "nodes"]
    verbs      = ["get", "watch", "list"]
  }
  rule {
    api_groups = ["coordination.k8s.io"]
    resources  = ["leases"]
    verbs      = ["get", "create", "update"]
  }
}

resource "kubernetes_cluster_role_binding" "fleet_server" {
  metadata {
    name = "fleet-server"
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.fleet_server.metadata[0].name
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.fleet_server.metadata[0].name
    namespace = kubernetes_namespace.elastic.metadata[0].name
  }
}

# Fleet Server Agent
resource "kubernetes_manifest" "fleet_server" {
  manifest = {
    apiVersion = "agent.k8s.elastic.co/v1alpha1"
    kind       = "Agent"
    metadata = {
      name      = "fleet-server"
      namespace = kubernetes_namespace.elastic.metadata[0].name
    }
    spec = yamldecode(file("${path.module}/k8s/fleet-server.yaml")).spec
  }
  
  depends_on = [
    kubernetes_manifest.kibana_cr,
    kubernetes_service_account.fleet_server,
    kubernetes_cluster_role_binding.fleet_server
  ]
}

# Service Account for Elastic Agent
resource "kubernetes_service_account" "elastic_agent" {
  metadata {
    name      = "elastic-agent"
    namespace = kubernetes_namespace.elastic.metadata[0].name
  }
}

resource "kubernetes_cluster_role" "elastic_agent" {
  metadata {
    name = "elastic-agent"
  }

  rule {
    api_groups = [""]
    resources  = ["pods", "nodes", "namespaces", "events", "services", "configmaps"]
    verbs      = ["get", "watch", "list"]
  }
  rule {
    api_groups = ["coordination.k8s.io"]
    resources  = ["leases"]
    verbs      = ["get", "create", "update"]
  }
  rule {
    api_groups = ["apps"]
    resources  = ["daemonsets", "deployments", "replicasets", "statefulsets"]
    verbs      = ["get", "watch", "list"]
  }
  rule {
    api_groups = ["batch"]
    resources  = ["jobs", "cronjobs"]
    verbs      = ["get", "watch", "list"]
  }
}

resource "kubernetes_cluster_role_binding" "elastic_agent" {
  metadata {
    name = "elastic-agent"
  }
  role_ref {
    api_group = "rbac.authorization.k8s.io"
    kind      = "ClusterRole"
    name      = kubernetes_cluster_role.elastic_agent.metadata[0].name
  }
  subject {
    kind      = "ServiceAccount"
    name      = kubernetes_service_account.elastic_agent.metadata[0].name
    namespace = kubernetes_namespace.elastic.metadata[0].name
  }
}

# Elastic Agent (DaemonSet for logs and APM)
resource "kubernetes_manifest" "elastic_agent" {
  manifest = {
    apiVersion = "agent.k8s.elastic.co/v1alpha1"
    kind       = "Agent"
    metadata = {
      name      = "elastic-agent"
      namespace = kubernetes_namespace.elastic.metadata[0].name
    }
    spec = yamldecode(file("${path.module}/k8s/elastic-agent.yaml")).spec
  }
  
  depends_on = [
    kubernetes_manifest.fleet_server,
    kubernetes_service_account.elastic_agent,
    kubernetes_cluster_role_binding.elastic_agent
  ]
}

# APM Service to expose elastic-agent APM server
resource "kubernetes_service" "apm" {
  metadata {
    name      = "apm"
    namespace = kubernetes_namespace.elastic.metadata[0].name
    labels = {
      app = "elastic-agent"
    }
  }

  spec {
    selector = {
      "agent.k8s.elastic.co/name" = "elastic-agent"
    }

    port {
      port        = 8200
      target_port = 8200
      protocol    = "TCP"
      name        = "apm"
    }

    type = "ClusterIP"
  }

  depends_on = [kubernetes_manifest.elastic_agent]
}
