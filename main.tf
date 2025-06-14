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
  }
}

module "nginx_conf" {
  source = "./modules/nginx_config"
}
