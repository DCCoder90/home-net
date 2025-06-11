resource "technitium_record" "ipv4_dns_record" {
  count = var.enable_ipv4 ? 1 : 0
  #zone       = var.zone_name
  domain     = var.domain_name
  type       = "A"
  ip_address = var.internal_host_ipv4
  ttl        = var.ttl
}

resource "technitium_record" "ipv6_dns_record" {
  count = var.enable_ipv6 ? 1 : 0
  #zone       = var.zone_name
  domain     = var.domain_name
  type       = "AAAA"
  ip_address = var.internal_host_ipv6
  ttl        = var.ttl
}