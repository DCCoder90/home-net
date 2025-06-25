locals {
  volumes_to_create = {
    for idx, vol_config in var.container_volumes : idx => vol_config
    if vol_config.volume_name != null && coalesce(vol_config.manage_volume_lifecycle, true)
  }

  all_labels = toset(concat(
    tolist(coalesce(var.labels, [])),
    var.icon != null ? [{
      label = "net.unraid.docker.icon",
      value = var.icon
    }] : [],
    var.web_ui != null ? [{
      label = "net.unraid.docker.webui",
      value = var.web_ui
    }] : []
  ))

  # Create a map of networks, keyed by network name, for easy lookups.
  # This resolves the "Invalid object key" error by using `net.name` (a string) as the key.
  all_networks_map = { for net in var.networks : net.name => net }

  # `network_mode` should be the first network from the list that does NOT have a static IP.
  potential_primary_networks = [
    for k, v in local.all_networks_map : k if v.ip_address == null
  ]

  # Define the effective network mode.
  effective_network_mode = (
    # Respect user's explicit choice of 'host' or 'none'.
    var.container_network_mode == "host" ? "host" :
    var.container_network_mode == "none" ? "none" :
    # If we have custom networks, pick a primary.
    length(local.all_networks_map) > 0 ? (
      # If we found a network without a static IP, use it as the primary.
      length(local.potential_primary_networks) > 0 ? local.potential_primary_networks[0] :
      # Otherwise  use 'bridge' as the primary
      # and attach all custom networks via networks_advanced.
      "bridge"
    ) :
    # If no custom networks are specified at all, default to 'bridge'.
    "bridge"
  )

  # Define the networks for the `networks_advanced` block
  advanced_network_attachments = (
    # If the effective network mode is 'host' or 'none', we don't need advanced networks.
    local.effective_network_mode == "host" || local.effective_network_mode == "none" ? {} :
    {
      for k, v in local.all_networks_map : k => v
    }
  )
}