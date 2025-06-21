terraform {
  required_providers {
    authentik = {
      source  = "goauthentik/authentik"
      version = "2025.4.0"
    }
  }
}

/*

No easy way to do this, have to manually assign applications to the outpost.  
See: https://github.com/goauthentik/terraform-provider-authentik/issues/310 
resource "authentik_outpost" "outpost" {
  name = "tf_${var.name}_outpost"
  protocol_providers = [
    authentik_provider_proxy.name.id
  ]
}
*/