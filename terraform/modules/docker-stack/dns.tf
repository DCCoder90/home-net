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
  
  # Really don't like having the ACLs hardcoded here...
  access_list_id           = each.value.network.internal ? local.npm_access_lists_by_name["Internal Only"] : local.npm_access_lists_by_name["Cloudflare"]
  internal_host_ipv4       = var.system.proxy_ip
  # If not using proxy auth, point to the service's static IP. If no static IP,
  # fall back to the service name, which NPM can use as a hostname.
  service_ipv4             = each.value.auth.proxy ? var.system.authentik.ip_address : coalesce(local.service_ip_addresses[each.key], each.value.service_name)
  admin_email              = var.system.network_admin_email
  dns_cloudflare_api_token = var.system.cloudflare_api_token
  external_host_ipv4       = var.system.public_facing_ip
}