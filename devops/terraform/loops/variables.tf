variable "environments" {
  description = "Map of environments and metadata."
  type = map(object({
    region  = string
    enabled = bool
    owners  = list(string)
  }))
  default = {
    dev = {
      region  = "us-east-1"
      enabled = true
      owners  = ["platform", "qa"]
    }
    qa = {
      region  = "us-east-2"
      enabled = true
      owners  = ["platform"]
    }
    prod = {
      region  = "sa-east-1"
      enabled = false
      owners  = ["platform", "security"]
    }
  }
}

variable "usernames" {
  description = "List used by for_each(list)."
  type        = list(string)
  default     = ["alice", "bob", "carol"]
}

variable "services" {
  description = "List of services used in nested iteration and filtering."
  type = list(object({
    name    = string
    port    = number
    tier    = string
    enabled = bool
  }))
  default = [
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
}

variable "worker_count" {
  description = "Number of instances created with count."
  type        = number
  default     = 3
}

variable "worker_prefix" {
  description = "Name prefix for resources created with count."
  type        = string
  default     = "worker"
}

variable "docker_host" {
  description = "Docker host socket/endpoint. If null, auto-detect Colima socket then fallback to default docker.sock."
  type        = string
  default     = null
  nullable    = true
}
