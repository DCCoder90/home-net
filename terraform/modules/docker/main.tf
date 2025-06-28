terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.6.2"
    }
  }
}

resource "docker_image" "main" {
  name = var.container_image
}

resource "docker_container" "container" {
  lifecycle {
    ignore_changes = [
      # Ignore changes to log_opts since unraid manages that
      log_opts,
    ]
  }
  name         = var.container_name
  image        = docker_image.main.image_id
  env          = var.environment_vars
  network_mode = local.effective_network_mode
  user         = var.container_user
  restart      = var.container_restart
  dns          = var.container_dns_servers
  privileged   = var.container_privileged_mode
  command      = var.commands
  runtime      = var.enable_gpu ? "nvidia" : null

  dynamic "networks_advanced" {
    for_each = local.advanced_network_attachments
    content {
      name         = networks_advanced.value.name
      ipv4_address = networks_advanced.value.ip_address
    }
  }

  dynamic "ports" {
    for_each = local.effective_network_mode != "host" ? var.container_ports : []
    content {
      internal = ports.value.internal
      external = ports.value.external # If null, Docker will assign a random host port
      ip       = ports.value.ip       # If null, Docker defaults to 0.0.0.0
      protocol = ports.value.protocol # If null, Docker defaults to tcp
    }
  }

  dynamic "mounts" {
    for_each = toset(coalesce(var.mounts, []))
    iterator = mount_iterator

    content {
      target    = split(":", mount_iterator.value)[1]
      source    = split(":", mount_iterator.value)[0]
      read_only = strcontains(mount_iterator.value, ":ro") ? true : false
      type      = "bind"
    }
  }

  dynamic "volumes" {
    for_each = var.container_volumes
    iterator = vol_iterator
    content {
      container_path = vol_iterator.value.container_path
      read_only      = vol_iterator.value.read_only
      host_path      = vol_iterator.value.host_path
      volume_name = vol_iterator.value.volume_name != null ? (
        coalesce(vol_iterator.value.manage_volume_lifecycle, true) ?
        docker_volume.managed_volumes[vol_iterator.key].name : # References module-created volume
        vol_iterator.value.volume_name                         # Uses pre-existing volume name directly
      ) : null
    }
  }

  //Capabilities causes a perpetual diff if it's empty, so we'll wrap it in a dynamic block
  dynamic "capabilities" {
    for_each = length(coalesce(var.container_capabilities.add, [])) > 0 || length(coalesce(var.container_capabilities.drop, [])) > 0 ? [1] : []
    content {
      add  = var.container_capabilities.add
      drop = var.container_capabilities.drop
    }
  }

  dynamic "labels" {
    for_each = local.all_labels

    content {
      label = labels.value.label
      value = labels.value.value
    }
  }
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