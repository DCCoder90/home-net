locals {
  # Determine the primary IP address for the service from its network configuration.
  service_ip_address = try(
    try([for n in var.service.network.networks : n.ip_address if n.name == "br1" && n.ip_address != null][0], null),
    try([for n in var.service.network.networks : n.ip_address if n.name == "br0" && n.ip_address != null][0], null),
    try([for n in var.service.network.networks : n.ip_address if n.ip_address != null][0], null)
  )

  # Generate environment variables for OAuth if it's enabled for the service.
  oauth_envs = var.service.auth.oauth.enabled ? [
    for env_name, output_key in var.service.auth.oauth.keys :
    format(
      "%s=%s",
      env_name,
      {
        "client_id"      = module.oauth_authentication[0].client_id,
        "client_secret"  = module.oauth_authentication[0].client_secret,
        "well_known_url" = module.oauth_authentication[0].provider_info_url
      }[output_key]
    )
  ] : []

  # Create a map of Nginx Proxy Manager access lists by name for easy lookup.
  npm_access_lists_by_name = { for al in data.nginxproxymanager_access_lists.access_lists.access_lists : al.name => al.id }

  domain = var.service.dns.domain_name == null ? "test.example" : var.service.dns.domain_name
  domain_name_parts = split(".", local.domain)
  zone_name = join(".", slice(local.domain_name_parts, length(local.domain_name_parts) - 2, length(local.domain_name_parts)))

}
