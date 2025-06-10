terraform {
  required_providers {
    authentik = {
      source  = "goauthentik/authentik"
      version = "2025.4.0"
    }
  }
}

provider "authentik" {
}

# SEE: https://registry.terraform.io/providers/goauthentik/authentik/latest/docs