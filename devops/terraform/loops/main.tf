
resource "docker_network" "services" {
  name   = "terraform-poc-network"
  driver = "bridge"
}

resource "docker_image" "service_image" {
  name         = "nginx:latest"
  keep_locally = true
}

resource "docker_container" "service" {
  for_each = local.service_map

  name  = "service-${each.key}"
  image = docker_image.service_image.image_id

  ports {
    internal = 80
    external = each.value.port
  }

  networks_advanced {
    name = docker_network.services.name
  }

  env = [
    "SERVICE_NAME=${each.value.name}",
    "SERVICE_PORT=${each.value.port}",
    "SERVICE_TIER=${each.value.tier}"
  ]

  labels {
    label = "service"
    value = each.key
  }

  labels {
    label = "tier"
    value = each.value.tier
  }

  labels {
    label = "managed_by"
    value = "terraform"
  }
}

resource "docker_container" "worker" {
  count = var.worker_count

  name  = "${var.worker_prefix}-${format("%02d", count.index + 1)}"
  image = docker_image.service_image.image_id

  networks_advanced {
    name = docker_network.services.name
  }

  env = [
    "WORKER_INDEX=${count.index + 1}",
    "WORKER_NAME=${var.worker_prefix}-${format("%02d", count.index + 1)}"
  ]

  labels {
    label = "type"
    value = "worker"
  }

  labels {
    label = "index"
    value = tostring(count.index + 1)
  }

  labels {
    label = "managed_by"
    value = "terraform"
  }
}

resource "terraform_data" "environment" {
  for_each = var.environments

  input = {
    name    = each.key
    region  = each.value.region
    enabled = each.value.enabled
    owners  = each.value.owners
  }
}

resource "terraform_data" "user" {
  for_each = toset(var.usernames)

  input = {
    username = each.value
    label    = "user-${each.value}"
  }
}

resource "terraform_data" "worker" {
  count = var.worker_count

  input = {
    index = count.index + 1
    name  = format("%s-%02d", var.worker_prefix, count.index + 1)
  }
}

resource "terraform_data" "env_service" {
  for_each = local.env_service_map

  input = each.value
}
