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
      version = "1.2.0"
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

    infisical = {
      source  = "infisical/infisical"
      version = "0.15.19"
    }
  }
}

module "nginx_conf" {
  source = "git@github.com:DCCoder90/home-tf-modules.git//nginx_config?ref=1.0.0"
}
