module "custom_network" {
  count    = length(local.creatable_networks) > 0 ? 1 : 0
  source   = "../../modules/docker-network"
  networks = local.creatable_networks
}

module "service_container" {
  source = "../../modules/docker"

  icon = var.service.icon
  web_ui = (
    var.service.network != null && var.service.network.ip_address != null && var.service.network.service_port != null ?
    "http://${var.service.network.ip_address}:${var.service.network.service_port}" : null
  )
  container_name         = var.service.service_name
  container_image        = var.service.image_name
  container_network_mode = var.service.network_mode
  environment_vars = toset(
    concat(coalesce(var.service.env, []),
      coalesce(local.oauth_envs[var.service.service_name], []),
  ))
  mounts                 = var.service.mounts
  container_capabilities = var.service.capabilities
  commands               = var.service.commands

  # Attach the container to custom networks defined in the stack, but only if the service
  # explicitly lists that network in its own configuration.
  networks = coalesce(var.service.network.networks, [])
}
