# Create a proxy provider

data "authentik_flow" "default-authorization-flow" {
  slug = "default-provider-authorization-implicit-consent"
}

resource "authentik_provider_proxy" "name" {
  name               = lower(var.Application_Name)
  internal_host      = var.internal_host
  external_host      = var.external_host
  basic_auth_password_attribute= var.Password_Attribute
  basic_auth_username_attribute = var.Username_Attribute
  basic_auth_enabled = true
  authorization_flow = data.authentik_flow.default-authorization-flow.id
}

resource "authentik_application" "name" {
  name              = var.Application_Name
  slug              = lower(replace(var.Application_Name," ", "-"))
  protocol_provider = authentik_provider_proxy.name.id
  group = var.Application_Group != "" ? var.Application_Group : null
  meta_description = var.Application_Description != "" ? var.Application_Description : null
}