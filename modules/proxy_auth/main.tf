terraform {
  required_providers {
    authentik = {
      source  = "goauthentik/authentik"
      version = "2025.4.0"
    }
  }
}

resource "authentik_outpost" "outpost" {
  name = "tf_${var.name}_outpost"
  protocol_providers = [
    authentik_provider_proxy.name.id
  ]
}