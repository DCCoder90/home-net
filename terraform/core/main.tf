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
    infisical = {
      source  = "infisical/infisical"
      version = "0.15.19"
    }
  }
}

locals {
  system = yamldecode(file("${path.module}/../../config/system.yaml"))

  # Load all core service definitions from config/core/*.yaml
  core_configs = merge([
    for f in fileset("${path.module}/../../config/core", "*.yaml") :
    yamldecode(file("${path.module}/../../config/core/${f}"))
  ]...)
}
