resource "kubernetes_config_map" "nginx_custom" {
  metadata {
    name      = "nginx-custom-config"
    namespace = kubernetes_namespace.app.metadata[0].name
  }
  data = {
    "nginx.conf" = <<-EOT
      events {}
      http {
        server {
          listen 80;
          location / {
            return 200 'Hello from custom NGINX!';
            add_header Content-Type text/plain;
          }
        }
      }
    EOT
  }
}

resource "kubernetes_deployment" "nginx_custom" {
  metadata {
    name      = "nginx-custom"
    namespace = kubernetes_namespace.app.metadata[0].name
    labels = {
      app = "nginx-custom"
    }
  }
  spec {
    replicas = 1
    selector {
      match_labels = {
        app = "nginx-custom"
      }
    }
    template {
      metadata {
        labels = {
          app = "nginx-custom"
        }
      }
      spec {
        container {
          name  = "nginx"
          image = "nginx:alpine"
          volume_mount {
            name       = "nginx-config"
            mount_path = "/etc/nginx/nginx.conf"
            sub_path   = "nginx.conf"
          }
          port {
            container_port = 80
          }
        }
        volume {
          name = "nginx-config"
          config_map {
            name = kubernetes_config_map.nginx_custom.metadata[0].name
          }
        }
      }
    }
  }
}

resource "kubernetes_service" "nginx_custom" {
  metadata {
    name      = "demo"
    namespace = kubernetes_namespace.app.metadata[0].name
  }
  spec {
    selector = {
      app = "nginx-custom"
    }
    port {
      port        = 80
      target_port = 80
    }
    type = "ClusterIP"
  }
}



resource "kubernetes_ingress_v1" "demo_ingress" {
  metadata {
    name = "demo-ingress"
    annotations = {
      "nginx.ingress.kubernetes.io/rewrite-target" = "/"
    }
    namespace = kubernetes_namespace.ingress_nginx.metadata[0].name
  }

  spec {
    rule {
      host = "demo.local"

      http {
        path {
          path     = "/"
          path_type = "Prefix"

          backend {
            service {
              name = kubernetes_service.nginx_custom.metadata[0].name
              port {
                number = 80
              }
            }
          }
        }
      }
    }
  }
}