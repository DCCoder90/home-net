terraform {
  cloud {
    organization = "DCCoders"
    workspaces {
      name = "home-core"
    }
  }

  required_providers {
    docker = {
      source  = "kreuzwerker/docker"
      version = "3.6.2"
    }

    nginxproxymanager = {
      source  = "Sander0542/nginxproxymanager"
      version = "1.2.0"
    }
  }
}

provider "docker" {
  host = "unix:///var/run/docker.sock"
}

provider "nginxproxymanager" {
  url      = var.nginx_proxy_address
  username = var.nginx_proxy_user
  password = var.nginx_proxy_pass
}

# Core Variables Loading
locals {
  system = yamldecode(file("${path.module}/../../config/system.yaml"))
}

data "docker_network" "br0" {
  name = "br0"
}

data "docker_network" "br1" {
  name = "br1"
}
