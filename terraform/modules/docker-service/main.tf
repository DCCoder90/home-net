

module "service_container" {
  source   = "../../modules/docker"

  icon                   = var.service.icon
  web_ui                 = try(var.service.network.service_port, null) != null && local.service_ip_address != null ? "http://${local.service_ip_address}:${var.service.network.service_port}" : null
  container_name         = var.service.service_name
  container_image        = var.service.image_name
  container_network_mode = var.service.network_mode
  enable_gpu             = var.service.enable_gpu
  environment_vars       = toset(concat(local.secret_envs, coalesce(local.oauth_envs, []), coalesce(var.service.env, [])))
  mounts                 = var.service.mounts
  container_capabilities = var.service.capabilities
  commands               = var.service.commands

  # Attach the container to custom networks defined in the stack, but only if the service
  # explicitly lists that network in its own configuration.
  # The docker module expects a list of objects with `name` and `ipv4_address`.
  networks = coalesce(var.service.network.networks, [])
}
