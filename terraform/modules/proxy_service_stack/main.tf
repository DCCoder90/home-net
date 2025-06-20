module "service_docker" {
  for_each = var.services
  source   = "../docker"

  container_name   = each.value.service_name
  container_image  = each.value.image_name
  attach_to_br0    = false
  attach_to_br1    = true
  br1_ipv4_addr    = each.value.ip_address
  environment_vars = each.value.env
}

module "service_dns" {
  source   = "../dns"
  for_each = var.services

  internal_only            = true
  service_port             = var.authentik_port
  zone_name                = var.zone_name
  domain_name              = each.value.domain_name
  access_list_id           = var.access_list_id
  internal_host_ipv4       = var.proxy_ip
  service_ipv4             = var.authentik_ip
  admin_email              = var.admin_email
  dns_cloudflare_api_token = var.cloudflare_token
  external_host_ipv4       = var.public_facing_ip
}

//Note: This will still require manually navigating to the service and setting up the credentials on the first run

module "authentication" {
  source   = "../proxy_auth"
  for_each = var.services

  internal_host               = "http://${each.value.ip_address}:${each.value.service_port}"
  external_host               = each.value.domain_name
  name                        = each.value.service_name
  username_attribute          = "${each.value.service_name}_username"
  password_attribute          = "${each.value.service_name}_password"
  create_access_group         = true
  access_group_name           = "tf_${each.value.service_name}"
  user_to_add_to_access_group = var.admin_username
  access_group_attributes = jsonencode(
    {
      "${each.value.service_name}_username" : each.value.username,
      "${each.value.service_name}_password" : each.value.password
    }
  )
}