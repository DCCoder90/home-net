terraform {
  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.6.2"
    }
  }
}

resource "docker_network" "private_network" {
  for_each   = var.networks
  name       = each.key
  internal   = each.value.internal
  driver     = each.value.driver
  attachable = true
  options    = each.value.options
}