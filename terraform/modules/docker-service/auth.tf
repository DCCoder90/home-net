

module "proxy_authentication" {
  source = "../proxy_auth"
  count = var.service.auth.enabled && var.service.auth.proxy.enabled ? 1 : 0

  group       = var.service.auth.group
  description = var.service.description
  # If a static IP is defined, use it. Otherwise, fall back to the service name,
  # which is resolvable within a Docker network.
  internal_host               = "http://${coalesce(local.service_ip_address, var.service.service_name)}:${var.service.network.service_port}"
  external_host               = var.service.dns.domain_name
  name                        = var.service.service_name
  username_attribute          = "${var.service.service_name}_username"
  password_attribute          = "${var.service.service_name}_password"
  create_access_group         = true
  access_group_name           = "tf_${var.service.service_name}"
  user_to_add_to_access_group = var.system.network_admin_username
  access_group_attributes = jsonencode(
    {
      "${var.service.service_name}_username" : data.infisical_secrets.secrets[0].secrets[coalesce(try(var.service.auth.proxy.user_secret, null), "${var.service.service_name}_username")].value,
      "${var.service.service_name}_password" : data.infisical_secrets.secrets[0].secrets[coalesce(try(var.service.auth.proxy.pass_secret, null), "${var.service.service_name}_password")].value
    }
  )
}

module "oauth_authentication" {
  source = "../oauth_auth"
  count = var.service.auth.enabled && var.service.auth.oauth.enabled ? 1 : 0

  group                       = var.service.auth.group
  description                 = var.service.description
  name                        = var.service.service_name
  create_access_group         = true
  access_group_name           = "tf_${var.service.service_name}"   
  user_to_add_to_access_group = var.system.network_admin_username 
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
