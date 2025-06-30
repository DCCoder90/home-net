locals {
  service_ip_addresses = {
    for service_key, service_config in var.stack.services : service_key => try(
      [for n in try(service_config.network.networks, []) : n.ip_address if n.name == "br1" && n.ip_address != null][0],
      [for n in try(service_config.network.networks, []) : n.ip_address if n.name == "br0" && n.ip_address != null][0],
      [for n in try(service_config.network.networks, []) : n.ip_address if n.ip_address != null][0],
      null
    )
  }

  # Networks that this module should actually create.
  # Filter out pre-existing host networks.
  creatable_networks = {
    for name, config in coalesce(var.stack.networks, {}) : name => config
    if !contains(var.system.existing_networks, name)
  }

 // services is a map of service configurations with environment variables processed
  // to replace secret placeholders with their actual values from your data source.
  services = {
    // Iterate over each service defined in your stack configuration.
    // 'service_name' is the key (e.g., "my-api"), and 'service_config' is the service's attribute map.
    for service_name, service_config in var.stack.services.definitions :
    service_name => merge(
      // We start by copying all original attributes from the service configuration.
      service_config,
      {
        // Then, we override the 'environment' attribute with the processed values.
        env = {
          // Iterate over each environment variable for the current service.
          // 'coalesce' ensures that if 'environment' is null, we safely iterate over an empty map.
          for env_name, env_value in coalesce(service_config.env, {}) :
          // The environment variable name (key) remains the same.
          env_name => reduce(
            // We iterate over the list of secret names that are candidates for replacement.
            var.stack.services.secrets,
            // The initial value for the reduction is the original environment variable's value.
            env_value,
            // This reducer function is applied for each secret name.
            // 'current_value' is the accumulated string, and 'secret_name' is the current secret being processed.
            (current_value, secret_name) -> replace(
              current_value,
              // This is the placeholder pattern to find, e.g., "${PATH_ROOT}".
              // The "$${" escapes Terraform's interpolation, so it looks for the literal characters.
              "$${${secret_name}}",
              // This is the value to substitute, fetched from your data source.
              // Note: This will raise an error if the secret_name does not exist in the data source,
              // which helps ensure all required secrets are present.
              data.infisical_secrets.generated_secrets.secrets[secret_name].value
            )
          )
        }
      }
    )
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