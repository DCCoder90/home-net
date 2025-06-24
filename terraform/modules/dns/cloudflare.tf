data "cloudflare_zone" "zone" {
  filter = {
    name = var.zone_name
  }
}

resource "cloudflare_dns_record" "ipv4_dns_record" {
  count   = var.internal_only ? 0 : var.enable_ipv4 ? 1 : 0
  zone_id = data.cloudflare_zone.zone.zone_id
  comment = "Managed via terraform"
  content = var.external_host_ipv4
  name    = var.domain_name
  proxied = var.proxied_domain
  ttl     = var.proxied_domain == true ? 1 : var.ttl
  type    = "A"
}

resource "cloudflare_dns_record" "ipv6_dns_record" {
  count   = var.internal_only ? 0 : var.enable_ipv6 ? 1 : 0
  zone_id = data.cloudflare_zone.zone.zone_id
  comment = "Managed via terraform"
  content = var.external_host_ipv6
  name    = var.domain_name
  proxied = var.proxied_domain
  ttl     = var.proxied_domain == true ? 1 : var.ttl
  type    = "AAAA"
}