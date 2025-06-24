locals {
  service_ip_addresses = {
    for service_key, service_config in var.stack.services : service_key => try(
      [for n in try(service_config.network.networks, []) : n.ipv4_address if n.name == "br1" && n.ipv4_address != null][0],
      [for n in try(service_config.network.networks, []) : n.ipv4_address if n.name == "br0" && n.ipv4_address != null][0],
      [for n in try(service_config.network.networks, []) : n.ipv4_address if n.ipv4_address != null][0],
      null
    )
  }

  # Networks that this module should actually create.
  # Filter out pre-existing host networks.
  creatable_networks = {
    for name, config in coalesce(var.stack.networks, {}) : name => config
    if !contains(var.system.existing_networks, name)
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

  oauth_envs = {
    for service_key, service_config in var.stack.services : service_key => (
      lookup(lookup(service_config, "auth", {}), "oauth", { enabled = false }).enabled ? [
        for env_name, output_key in lookup(lookup(service_config, "auth", {}), "oauth", { keys = {} }).keys :
        format(
          "%s=%s",
          env_name,
          # This map translates the key from your YAML (e.g., "client_id")
          # to the corresponding attribute of the created Authentik resource.
          {
            "client_id" : module.oauth_authentication[service_key].client_id,
            "client_secret" : module.oauth_authentication[service_key].client_secret,
            "well_known_url" : module.oauth_authentication[service_key].provider_info_url
          }[output_key] # Use the value from the YAML keys map to look up the correct output
        )
      ] : []
    )
  }

  npm_access_lists_by_name = { for al in data.nginxproxymanager_access_lists.access_lists.access_lists : al.name => al.id }
}