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

module "service_container" {
  for_each = var.stack.services
  source   = "../../modules/docker"

  icon = each.value.icon
  web_ui = try(each.value.network.service_port, null) != null && local.service_ip_addresses[each.key] != null ? "http://${local.service_ip_addresses[each.key]}:${each.value.network.service_port}" : null
  container_name         = each.value.service_name
  container_image        = each.value.image_name
  container_network_mode = each.value.network_mode
  environment_vars       = toset(concat(coalesce(var.stack.env, []), local.processed_envs[each.key], coalesce(var.stack.env, []), coalesce(local.oauth_envs[each.key], [])))
  mounts                 = concat(coalesce(var.stack.mounts, []), coalesce(each.value.mounts, []))
  container_capabilities = each.value.capabilities
  commands               = each.value.commands

  # Attach the container to custom networks defined in the stack, but only if the service
  # explicitly lists that network in its own configuration.
  # The docker module expects a list of objects with `name` and `ipv4_address`.
  networks = coalesce(each.value.network.networks, [])
}
