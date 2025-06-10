terraform {
  required_providers {
    authentik = {
      source  = "goauthentik/authentik"
      version = "2025.4.0"
    }

    nginxproxymanager = {
      source = "Sander0542/nginxproxymanager"
      version = "1.1.1"
    }

    portainer = {
      source = "portainer/portainer"
      version = "1.4.3"
    }
  }
}

provider "authentik" {
}

# SEE: https://registry.terraform.io/providers/goauthentik/authentik/latest/docs
# SEE: https://registry.terraform.io/providers/Sander0542/nginxproxymanager/latest
# SEE: https://registry.terraform.io/providers/portainer/portainer/latest