terraform {
  required_providers {
    random = {
      source = "hashicorp/random"
    }
    nginxproxymanager = {
      source  = "Sander0542/nginxproxymanager"
      version = "1.2.0"
    }
    infisical = {
      source  = "infisical/infisical"
      version = "0.15.19"
    }
  }
}