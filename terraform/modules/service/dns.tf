data "nginxproxymanager_access_lists" "access_lists" {}

module "service_dns" {
  source = "../dns"
  count  = var.service.dns.enabled ? 1 : 0

  internal_only            = var.service.network.internal
  service_port             = var.service.auth.proxy ? var.system.authentik.port : var.service.network.service_port
  zone_name                = var.stack.zone_name
  domain_name              = var.service.dns.domain_name
  access_list_id           = var.service.network.internal ? local.npm_access_lists_by_name["Internal Only"] : local.npm_access_lists_by_name["Cloudflare"]
  internal_host_ipv4       = var.system.proxy_ip
  service_ipv4             = var.service.auth.proxy ? var.system.authentik.ip_address : var.service.network.ip_address
  admin_email              = var.system.network_admin_email
  dns_cloudflare_api_token = var.system.cloudflare_api_token
  external_host_ipv4       = var.system.public_facing_ip
}