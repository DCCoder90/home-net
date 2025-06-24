
resource "random_password" "service_password" {
  for_each = {
    for k, v in var.stack.services : k => v if lookup(lookup(v, "auth", {}), "enabled", false) == true && lookup(lookup(v, "auth", {}), "proxy", false) == true
  }

  length           = 24
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}


module "proxy_authentication" {
  source = "../proxy_auth"
  for_each = {
    for k, v in var.stack.services : k => v if lookup(lookup(v, "auth", {}), "enabled", false) == true && lookup(lookup(v, "auth", {}), "proxy", false) == true
  }

  group                       = each.value.auth.group
  description                 = each.value.description
  internal_host               = "http://${local.service_ip_addresses[each.key]}:${each.value.network.service_port}"
  external_host               = each.value.dns.domain_name
  name                        = each.value.service_name
  username_attribute          = "${each.value.service_name}_username"
  password_attribute          = "${each.value.service_name}_password"
  create_access_group         = true
  access_group_name           = "tf_${each.value.service_name}"
  user_to_add_to_access_group = var.system.network_admin_username
  access_group_attributes = jsonencode(
    {
      "${each.value.service_name}_username" : "admin",
      "${each.value.service_name}_password" : random_password.service_password[each.key].result
    }
  )
}

module "oauth_authentication" {
  source = "../oauth_auth"
  for_each = {
    for k, v in var.stack.services : k => v if lookup(lookup(v, "auth", {}), "enabled", false) == true && lookup(lookup(lookup(v, "auth", {}), "oauth", {}), "enabled", false) == true
  }

  group                       = each.value.auth.group
  description                 = each.value.description
  name                        = each.value.service_name
  create_access_group         = true
  access_group_name           = "tf_${each.value.service_name}"   //Update to allow for using custom groups rather than just generated
  user_to_add_to_access_group = var.system.network_admin_username //Update to allow for a list of users
  allowed_redirect_uris = concat(
    [
      {
        matching_mode = "strict",
      url = "https://${each.value.dns.domain_name}" }
    ],
    [
      for uri_path in coalesce(each.value.auth.oauth.redirect_uris, []) : {
        matching_mode = "strict",
        url           = "https://${each.value.dns.domain_name}/${uri_path}"
      }
    ]
  )
}
