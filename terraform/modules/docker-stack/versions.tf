terraform {
  required_providers {
    random = {
      source = "hashicorp/random"
    }
    nginxproxymanager = {
      source  = "Sander0542/nginxproxymanager"
      version = "1.1.1"
    }
  }
}