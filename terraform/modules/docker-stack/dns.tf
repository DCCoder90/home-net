data "nginxproxymanager_access_lists" "access_lists" {}

module "service_dns" {
  source = "../dns"
  for_each = {
    for k, v in var.stack.services : k => v if lookup(v, "dns").enabled == true
  }

  internal_only            = each.value.network.internal
  service_port             = each.value.auth.proxy ? var.system.authentik.port : each.value.network.service_port
  zone_name                = var.stack.zone_name
  domain_name              = each.value.dns.domain_name
  access_list_id           = each.value.network.internal ? local.npm_access_lists_by_name["Internal Only"] : local.npm_access_lists_by_name["Cloudflare"]
  internal_host_ipv4       = var.system.proxy_ip
  service_ipv4             = each.value.auth.proxy ? var.system.authentik.ip_address : each.value.network.ip_address
  admin_email              = var.system.network_admin_email
  dns_cloudflare_api_token = var.system.cloudflare_api_token
  external_host_ipv4       = var.system.public_facing_ip
}