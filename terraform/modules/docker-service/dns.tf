data "nginxproxymanager_access_lists" "access_lists" {}

module "service_dns" {
  source = "../dns"
  count = var.service.dns.enabled ? 1 : 0

  internal_only = var.service.network.internal
  service_port  = var.service.auth.proxy.enabled ? var.system.authentik.port : var.service.network.service_port
  zone_name     = local.zone_name
  domain_name   = var.service.dns.domain_name

  # Really don't like having the ACLs hardcoded here...
  access_list_id     = var.service.network.internal ? local.npm_access_lists_by_name["Internal Only"] : local.npm_access_lists_by_name["Cloudflare"]
  internal_host_ipv4 = var.system.proxy_ip
  # If not using proxy auth, point to the service's static IP. If no static IP,
  # fall back to the service name, which NPM can use as a hostname.
  service_ipv4             = var.service.auth.proxy.enabled ? var.system.authentik.ip_address : coalesce(local.service_ip_address, var.service.service_name)
  admin_email              = var.system.network_admin_email
  dns_cloudflare_api_token = var.system.cloudflare_api_token
  external_host_ipv4       = var.system.public_facing_ip
}