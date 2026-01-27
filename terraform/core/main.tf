terraform {
  cloud {
    organization = "DCCoders"
    workspaces {
      name = "home-core"
    }
  }

  required_providers {
    nginxproxymanager = {
      source  = "Sander0542/nginxproxymanager"
      version = "1.2.0"
    }
    # Cloudflare/Infisical can go here if you move DNS zones to Core later
  }
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
