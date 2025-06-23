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

  icon                   = each.value.icon
  web_ui                 = (
    each.value.network != null && each.value.network.ip_address != null && each.value.network.service_port != null ?
    "http://${each.value.network.ip_address}:${each.value.network.service_port}" : null
  )
  container_name         = each.value.service_name
  container_image        = each.value.image_name
  attach_to_br0          = contains(coalesce(each.value.network.networks, []), "br0")
  attach_to_br1          = contains(coalesce(each.value.network.networks, []), "br1")
  br0_ipv4_addr          = contains(coalesce(each.value.network.networks, []), "br0") ? each.value.network.ip_address : null
  br1_ipv4_addr          = contains(coalesce(each.value.network.networks, []), "br1") ? each.value.network.ip_address : null
  environment_vars       = toset(concat(coalesce(var.stack.env, []), local.processed_envs[each.key], coalesce(var.stack.env, []), coalesce(local.oauth_envs[each.key], [])))
  mounts                 = concat(coalesce(var.stack.mounts, []), coalesce(each.value.mounts, []))
  container_capabilities = each.value.capabilities
  commands               = each.value.commands

  # Attach the container to custom networks defined in the stack, but only if the service
  # explicitly lists that network in its own configuration.
  networks = length(local.creatable_networks) > 0 ? [
    for net_name in coalesce(each.value.network.networks, []) : module.custom_network[0].networks[net_name].id if net_name != "br0" && net_name != "br1"
  ] : []
}
