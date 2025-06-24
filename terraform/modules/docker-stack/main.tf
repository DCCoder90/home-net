resource "random_password" "generated" {
  for_each = toset(coalesce(var.stack.generated_secrets, []))

  length  = 48
  special = true
}


module "custom_network" {
  count    = length(local.creatable_networks) > 0 ? 1 : 0
  source   = "../../modules/docker-network"
  networks = local.creatable_networks
}

module "service" {
  for_each = local.build_service
  source   = "../../modules/service"

  service = each.value
  system = var.system
  ignore_auth = true
  generated_secrets = local.generated_secrets
}

locals {
  build_service = {
    for service_name, service in var.stack.services : service_name => {
      service_name = service_name
      image_name   = service.image_name
      networks     = service.network.networks
      env          = concat(
        coalesce(var.stack.env, []),
        local.processed_envs[service_name],
        coalesce(var.stack.env, [])
      )
      mounts       = var.stack.mounts
      capabilities = service.capabilities
    }
  }
}