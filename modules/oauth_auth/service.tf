resource "random_string" "example" {
  length           = 16
  special          = false
  upper            = true
  lower            = true
  numeric          = true
  override_special = "!@#$%^&*"
}

data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

data "authentik_flow" "default-invalidation-flow" {
  slug = "default-invalidation-flow"
}

# Create an OAuth2 Provider

resource "authentik_provider_oauth2" "name" {
  name                  = lower(replace(var.name, " ", "-"))
  client_id             = random_string.example.result
  authorization_flow    = authentik_flow.default-authorization-flow.id
  invalidation_flow     = authentik_flow.default-invalidation-flow.id
  //allowed_redirect_uris = var.redirect_uris
}

resource "authentik_application" "name" {
  name              = var.name
  meta_description  = var.description
  group             = var.group
  slug              = lower(replace(var.name, " ", "-"))
  protocol_provider = authentik_provider_oauth2.name.id
}