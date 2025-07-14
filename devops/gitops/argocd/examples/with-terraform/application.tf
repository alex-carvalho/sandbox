resource "argocd_application" "application" {
  metadata {
    name      = "application-terraform"
    namespace = "argocd"
    labels = {
      using_sync_policy_options = "true"
    }
  }

  spec {
    destination {
      server    = "https://kubernetes.default.svc"
      namespace = "apps"
    }

    source {
      repo_url        = "https://github.com/alex-carvalho/sandbox.git"
      path            = "devops/gitops/argocd/examples/guestbook/manifests"
      target_revision = "HEAD"
    }
    sync_policy {
      automated {
        prune     = true
        self_heal = true
      }
    }
  }
}