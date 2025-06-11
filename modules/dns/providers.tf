terraform {
  required_providers {
    nginxproxymanager = {
      source  = "Sander0542/nginxproxymanager"
      version = "1.1.1"
    }

    cloudflare = {
      source  = "cloudflare/cloudflare"
      version = "~> 5"
    }

    technitium = {
      source  = "kevynb/technitium"
      version = "0.2.0"
    }
  }
}