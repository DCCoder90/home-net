locals {
  volumes_to_create = {
    for idx, vol_config in var.container_volumes : idx => vol_config
    if vol_config.volume_name != null && coalesce(vol_config.manage_volume_lifecycle, true)
  }
}