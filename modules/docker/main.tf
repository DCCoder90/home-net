terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.6.1"
    }
  }
}

resource "docker_container" "container" {
  name  = var.container_name
  image = docker_image.main.image_id
  env = var.environment_vars
  attach = false
  network_mode = var.container_network_mode

  dynamic "networks_advanced" {
    # Only attach to advanced networks if mode is not 'none' or 'host'
    for_each = var.attach_to_br1 && !contains(["none", "host"], var.container_network_mode) ? [1] : []
    content {
      name         = data.docker_network.main_host.name
      ipv4_address = var.br1_ipv4_addr 
    }
  }

  dynamic "networks_advanced" {
    for_each = var.attach_to_br0 && !contains(["none", "host"], var.container_network_mode) ? [1] : []
    content {
      name         = data.docker_network.secondary_host.name 
      ipv4_address = var.br0_ipv4_addr                     
    }
  }

  dynamic "ports" {
    # Only map ports if network mode is not 'none' or 'host'
    for_each = !contains(["none", "host"], var.container_network_mode) ? var.container_ports : []
    content {
      internal = ports.value.internal
      external = ports.value.external # If null, Docker will assign a random host port
      ip       = ports.value.ip       # If null, Docker defaults to 0.0.0.0
      protocol = ports.value.protocol # If null, Docker defaults to tcp
    }
  }

  dynamic "volumes" {
    for_each = var.container_volumes
    # Create a unique key for each volume entry to iterate over
    iterator = vol_iterator # Using an explicit iterator name for clarity
    content {
      container_path = vol_iterator.value.container_path
      read_only      = vol_iterator.value.read_only
      host_path      = vol_iterator.value.host_path                                                              # Will be null if volume_name is used
      volume_name    = vol_iterator.value.volume_name != null ? (
        coalesce(vol_iterator.value.manage_volume_lifecycle, true) ?
        docker_volume.managed_volumes[vol_iterator.key].name : # References module-created volume
        vol_iterator.value.volume_name                         # Uses pre-existing volume name directly
      ) : null
    }
  }

  user    = var.container_user
  restart = var.container_restart
  dns     = !contains(["none", "host"], var.container_network_mode) ? var.container_dns_servers : null # Only set DNS if mode is not 'none' or 'host'
  privileged = var.container_privileged_mode
}

resource "docker_image" "main" {
  name = var.container_image
}

resource "docker_volume" "managed_volumes" {
  for_each = local.volumes_to_create

  name        = each.value.volume_name
  driver      = each.value.driver
  driver_opts = each.value.driver_opts

  dynamic "labels" {
    for_each = each.value.labels == null ? {} : each.value.labels
    content {
      label = labels.key
      value = labels.value
    }
  }
}