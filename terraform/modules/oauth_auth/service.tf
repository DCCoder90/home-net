resource "random_string" "client_id" {
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

data "authentik_flow" "default-provider-authorization-implicit-consent" {
  slug = "default-provider-authorization-implicit-consent"
}

data "authentik_property_mapping_provider_scope" "scope-email" {
  name = "authentik default OAuth Mapping: OpenID 'email'"
}

data "authentik_property_mapping_provider_scope" "scope-profile" {
  name = "authentik default OAuth Mapping: OpenID 'profile'"
}

data "authentik_property_mapping_provider_scope" "scope-openid" {
  name = "authentik default OAuth Mapping: OpenID 'openid'"
}

# Create an OAuth2 Provider

# To get the the ID and other info about a certificate

data "authentik_certificate_key_pair" "generated" {
  name = "authentik Self-signed Certificate"
}

# Then use `data.authentik_certificate_key_pair.generated.id`

resource "authentik_provider_oauth2" "name" {
  name               = lower(replace(var.name, " ", "-"))
  client_id          = random_string.client_id.result
  authorization_flow = data.authentik_flow.default-authorization-flow.id
  invalidation_flow  = data.authentik_flow.default-invalidation-flow.id
  signing_key        = data.authentik_certificate_key_pair.generated.id

  allowed_redirect_uris = var.allowed_redirect_uris

  property_mappings = [
    data.authentik_property_mapping_provider_scope.scope-email.id,
    data.authentik_property_mapping_provider_scope.scope-profile.id,
    data.authentik_property_mapping_provider_scope.scope-openid.id,
  ]
}

resource "authentik_application" "name" {
  name              = var.name
  meta_description  = var.description
  group             = var.group
  slug              = lower(replace(var.name, " ", "-"))
  protocol_provider = authentik_provider_oauth2.name.id
}

data "authentik_provider_oauth2_config" "oauth_config" {
  name = authentik_provider_oauth2.name.name
}