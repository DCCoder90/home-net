resource "random_password" "service_password" {
  count = var.service.auth.proxy && var.service.auth.enabled ? 1 : 0

  length           = 24
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}


module "proxy_authentication" {
  source = "../proxy_auth"
  count  = var.service.auth.proxy && var.service.auth.enabled ? 1 : 0

  group       = var.service.auth.group
  description = var.service.description
  internal_host = "http://${try(
    [for n in var.service.network.networks : n.ip_address if n.name == "br1"][0],
    [for n in var.service.network.networks : n.ip_address if n.name == "br0"][0],
    var.service.network.networks[0].ip_address
  )}:${var.service.network.service_port}"
  external_host               = var.service.dns.domain_name
  name                        = var.service.service_name
  username_attribute          = "${var.service.service_name}_username"
  password_attribute          = "${var.service.service_name}_password"
  create_access_group         = true
  access_group_name           = "tf_${var.service.service_name}"
  user_to_add_to_access_group = var.system.network_admin_username
  access_group_attributes = jsonencode(
    {
      "${var.service.service_name}_username" : "admin",
      "${var.service.service_name}_password" : random_password.service_password[0].result
    }
  )
}

module "oauth_authentication" {
  source = "../oauth_auth"
  count  = var.service.auth.enabled && var.service.auth.oauth.enabled ? 1 : 0

  group                       = var.service.auth.group
  description                 = var.service.description
  name                        = var.service.service_name
  create_access_group         = true
  access_group_name           = "tf_${var.service.service_name}"  //Update to allow for using custom groups rather than just generated
  user_to_add_to_access_group = var.system.network_admin_username //Update to allow for a list of users
  allowed_redirect_uris = concat(
    [
      {
        matching_mode = "strict",
      url = "https://${var.service.dns.domain_name}" }
    ],
    [
      for uri_path in coalesce(var.service.auth.oauth.redirect_uris, []) : {
        matching_mode = "strict",
        url           = "https://${var.service.dns.domain_name}/${uri_path}"
      }
    ]
  )
}
