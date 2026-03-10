usernames = [
  "ana",
  "bruno",
  "carlos"
]

worker_count  = 2
worker_prefix = "batch"

services = [
  {
    name    = "api"
    port    = 8080
    tier    = "backend"
    enabled = true
  },
  {
    name    = "web"
    port    = 3000
    tier    = "frontend"
    enabled = true
  },
  {
    name    = "jobs"
    port    = 9090
    tier    = "worker"
    enabled = false
  }
]
