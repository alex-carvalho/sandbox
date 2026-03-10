output "docker_network" {
  description = "Docker network created for services."
  value = {
    name = docker_network.services.name
    id   = docker_network.services.id
  }
}

output "service_containers" {
  description = "Docker containers for each enabled service (for_each)."
  value = {
    for name, container in docker_container.service : name => {
      id             = container.id
      container_name = container.name
      image          = container.image
      ports          = [for p in container.ports : "${p.external}:${p.internal}"]
      env            = container.env
    }
  }
}

output "worker_containers" {
  description = "Worker containers created with count."
  value = [
    for worker in docker_container.worker : {
      id   = worker.id
      name = worker.name
      env  = worker.env
    }
  ]
}

output "environments_by_name" {
  description = "for_each(map) result."
  value = {
    for name, item in terraform_data.environment : name => item.output
  }
}

output "users" {
  description = "for_each(list) result."
  value = [
    for user in terraform_data.user : user.output
  ]
}

output "worker_names" {
  description = "count result."
  value = [
    for worker in terraform_data.worker : worker.output.name
  ]
}

output "service_names_enabled_only" {
  description = "for expression with filter (if condition)."
  value       = local.enabled_service_names
}

output "env_service_matrix" {
  description = "Nested loop transformation (environment x service)."
  value = {
    for key, item in terraform_data.env_service : key => item.output
  }
}
