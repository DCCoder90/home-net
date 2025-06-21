module "delugevpn_service" {
  source   = "../../modules/docker"

  icon = var.service.icon
  web_ui = "http://${var.service.ip_address}:${var.service.service_port}"
  container_name   = var.service.service_name
  container_image  = var.service.image_name
  attach_to_br0    = false
  attach_to_br1    = true
  br1_ipv4_addr    = var.service.ip_address
  environment_vars = var.service.env
  mounts           = var.service.mounts
  container_capabilities = var.service.capabilities
}

module "service_dns" {
  source   = "../dns"
  for_each = var.stack.services

  internal_only            = each.value.internal
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