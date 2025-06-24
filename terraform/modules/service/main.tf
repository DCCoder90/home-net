module "custom_network" {
  count    = length(local.creatable_networks) > 0 ? 1 : 0
  source   = "../../modules/docker-network"
  networks = local.creatable_networks
}

module "service_container" {
  source = "../../modules/docker"

  icon = var.service.value.icon
  web_ui = (
    var.service.value.network != null && var.service.value.network.ip_address != null && var.service.value.network.service_port != null ?
    "http://${each.value.network.ip_address}:${each.value.network.service_port}" : null
  )
  container_name         = var.service.value.service_name
  container_image        = var.service.value.image_name
  container_network_mode = var.service.value.network_mode
  environment_vars       = toset(concat(coalesce(var.stack.env, []), local.processed_envs[var.service.service_name], coalesce(var.stack.env, []), coalesce(local.oauth_envs[var.service.service_name], [])))
  mounts                 = concat(coalesce(var.stack.mounts, []), coalesce(var.service.value.mounts, []))
  container_capabilities = var.service.value.capabilities
  commands               = var.service.value.commands

  # Attach the container to custom networks defined in the stack, but only if the service
  # explicitly lists that network in its own configuration.
  networks = length(local.creatable_networks) > 0 ? [
    for net in coalesce(each.value.network.networks, []) : module.custom_network[0].networks[net.name].id
  ] : []
}
