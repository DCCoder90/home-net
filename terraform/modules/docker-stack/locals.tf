locals {
  # Networks that this module should actually create.
  # Filter out br0 and br1 as they are pre-existing host networks.
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

  npm_access_lists_by_name = { for al in data.nginxproxymanager_access_lists.access_lists.access_lists : al.name => al.id }
}