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
  # Merge static configuration from system.yaml with runtime variables.
  # This creates a single, comprehensive configuration object to pass to modules.
  system = merge(
    yamldecode(file("${path.module}/../config/system.yaml")),
    {
      nginx_proxy_address    = var.nginx_proxy_address
      nginx_proxy_user       = var.nginx_proxy_user
      nginx_proxy_pass       = var.nginx_proxy_pass
      cloudflare_api_token   = var.cloudflare_api_token
      network_admin_email    = var.network_admin_email
      technitium_api_token   = var.technitium_api_token
      technitium_host        = var.technitium_host
      public_facing_ip       = var.public_facing_ip
      network_admin_username = var.network_admin_username
      vpn_user               = var.vpn_user
      vpn_pass               = var.vpn_pass
    }
  )
  stacks = merge([
    for f in fileset("${path.module}/../config/stacks", "*.yaml") :
    yamldecode(file("${path.module}/../config/stacks/${f}"))
  ]...)
  services = merge([
    for f in fileset("${path.module}/../config/services", "*.yaml") :
    yamldecode(file("${path.module}/../config/services/${f}"))
  ]...)
}