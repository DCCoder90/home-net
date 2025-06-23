locals {
  volumes_to_create = {
    for idx, vol_config in var.container_volumes : idx => vol_config
    if vol_config.volume_name != null && coalesce(vol_config.manage_volume_lifecycle, true)
  }

  effective_network_mode = (
    var.container_network_mode == "host" ? "host" : 
    var.attach_to_br1 ? "br1" : 
    var.attach_to_br0 ? "br0" : 
    ( 
      var.container_network_mode == null ? ( 
        length(var.networks) > 0 ? element(var.networks, 0) : 
        "bridge" # Fallback to "bridge" if no custom networks are specified
      ) :
      var.container_network_mode # Fallback to the provided container_network_mode if it's not null
      )
    )

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
}