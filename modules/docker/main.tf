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
      # Ignore changes to tags, e.g. because a management agent
      # updates these based on some ruleset managed elsewhere.
      log_opts,
    ]
  }
  name  = var.container_name
  image = docker_image.main.image_id
  env = var.environment_vars
  network_mode = local.effective_network_mode
  user = var.container_user
  restart = var.container_restart
  dns     = local.effective_network_mode != "host" ? var.container_dns_servers : null
  privileged = var.container_privileged_mode

  dynamic "networks_advanced" {
    for_each = var.attach_to_br1 && local.effective_network_mode != "host" ? [1] : []
    content {
      name         = data.docker_network.main_host.id
      ipv4_address = var.br1_ipv4_addr 
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