data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

data "authentik_flow" "default-invalidation-flow" {
  slug = "default-invalidation-flow"
}

resource "authentik_provider_proxy" "name" {
  name                          = lower(replace(var.name, " ", "-"))
  internal_host                 = var.internal_host
  external_host                 = "http://${var.external_host}"
  basic_auth_password_attribute = var.password_attribute
  basic_auth_username_attribute = var.username_attribute
  basic_auth_enabled            = true
  authorization_flow            = data.authentik_flow.default-authorization-flow.id
  invalidation_flow             = data.authentik_flow.default-invalidation-flow.id
}

resource "authentik_application" "name" {
  name              = var.name
  meta_description  = var.description
  group             = var.group
  slug              = "tf_${lower(replace(var.name, " ", "-"))}"
  protocol_provider = authentik_provider_proxy.name.id
}