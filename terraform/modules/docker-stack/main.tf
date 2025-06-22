locals {
  # Networks that this module should actually create.
  # We filter out br0 and br1 as they are pre-existing host networks.
  creatable_networks = {
    for name, config in coalesce(var.stack.networks, {}) : name => config
    if name != "br0" && name != "br1"
  }

  # A map of generated secrets, with the secret name as the key.
  generated_secrets = {
    for key, secret in random_password.generated : key => secret.result
  }

  # For each service, process its environment variables to substitute secret placeholders.
  processed_envs = {
    for service_key, service_config in var.stack.services : service_key => [
      for env_string in coalesce(service_config.env, []) : (
        # This pattern ensures that for each env var, we get either the string with the secret
        # replaced, or the original string if no placeholder was found.
        # The one() function enforces a "one secret per line" rule, failing if multiple placeholders exist.
        one(concat(
          [
            for secret_name, secret_value in local.generated_secrets :
            replace(env_string, format("$${%s}", secret_name), secret_value)
            if strcontains(env_string, format("$${%s}", secret_name))
          ],
          # This list provides the fallback original string if no secret placeholder was found.
          !anytrue([for name in keys(local.generated_secrets) : strcontains(env_string, format("$${%s}", name))]) ? [env_string] : []
        ))
      )
    ]
  }
}

resource "random_password" "generated" {
  for_each = toset(coalesce(var.stack.generated_secrets, []))

  length  = 48
  special = true
}

module "karakeep_network" {
  count    = length(local.creatable_networks) > 0 ? 1 : 0
  source   = "../../modules/docker-network"
  networks = local.creatable_networks
}

module "service_container" {
  for_each = var.stack.services
  source   = "../../modules/docker"

  icon                   = each.value.icon
  web_ui                 = "http://${each.value.ip_address}:${each.value.service_port}"
  container_name         = each.value.service_name
  container_image        = each.value.image_name
  attach_to_br0          = contains(coalesce(each.value.networks, []), "br0")
  attach_to_br1          = contains(coalesce(each.value.networks, []), "br1")
  br0_ipv4_addr          = contains(coalesce(each.value.networks, []), "br0") ? each.value.ip_address : null
  br1_ipv4_addr          = contains(coalesce(each.value.networks, []), "br1") ? each.value.ip_address : null
  environment_vars       = toset(concat(coalesce(var.stack.env, []), local.processed_envs[each.key]))
  mounts                 = coalesce(each.value.mounts, [])
  container_capabilities = each.value.capabilities

  # Attach the container to custom networks defined in the stack, but only if the service
  # explicitly lists that network in its own configuration.
  networks = length(local.creatable_networks) > 0 ? [
    for net_name in coalesce(each.value.networks, []) : module.karakeep_network[0].networks[net_name].id if net_name != "br0" && net_name != "br1"
  ] : []
}