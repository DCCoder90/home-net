
resource "random_password" "service_password" {
  for_each = {
      for k, v in var.stack.services : k => v if lookup(v, "auth").enabled == true && lookup(v, "auth").proxy == true
    }

  length           = 24
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}


module "authentication" {
  source   = "../proxy_auth"
  for_each = {
      for k, v in var.stack.services : k => v if lookup(v, "auth").enabled == true && lookup(v, "auth").proxy == true
    }

  internal_host               = "http://${each.value.network.ip_address}:${each.value.network.service_port}"
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