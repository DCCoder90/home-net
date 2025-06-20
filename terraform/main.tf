terraform {
  cloud {

    organization = "DCCoders"

    workspaces {
      name = "home"
    }
  }

  required_providers {
    nginxproxymanager = {
      source  = "Sander0542/nginxproxymanager"
      version = "1.1.1"
    }

    docker = {
      source  = "kreuzwerker/docker"
      version = "3.6.2"
    }

    technitium = {
      source  = "kevynb/technitium"
      version = "0.2.0"
    }

    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5"
    }

    authentik = {
      source  = "goauthentik/authentik"
      version = "2025.4.0"
    }
  }
}

module "nginx_conf" {
  source = "./modules/nginx_config"
}

locals {
  system = yamldecode(file("${path.module}/../config/system.yaml"))
  stacks = yamldecode(file("${path.module}/../config/stacks.yaml"))
  services = yamldecode(file("${path.module}/../config/services.yaml"))
}