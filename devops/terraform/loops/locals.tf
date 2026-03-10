
locals {
  normalized_services = [
    for service in var.services : {
      name    = lower(service.name)
      port    = service.port
      tier    = lower(service.tier)
      enabled = service.enabled
    }
  ]

  enabled_service_names = [
    for service in local.normalized_services : service.name
    if service.enabled
  ]

  env_service_pairs = flatten([
    for env_name, env in var.environments : [
      for service in local.normalized_services : {
        key         = "${env_name}-${service.name}"
        environment = env_name
        region      = env.region
        owners      = env.owners
        service     = service.name
        port        = service.port
        tier        = service.tier
      } if env.enabled && service.enabled
    ]
  ])

  env_service_map = {
    for pair in local.env_service_pairs : pair.key => pair
  }

  service_map = {
    for service in local.normalized_services : service.name => service
    if service.enabled
  }
}