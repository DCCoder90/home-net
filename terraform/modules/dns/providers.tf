terraform {
  required_providers {
    nginxproxymanager = {
      source  = "Sander0542/nginxproxymanager"
      version = "1.2.0"
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