resource "nginxproxymanager_certificate_letsencrypt" "certificate" {
  domain_names = [var.domain_name]

  letsencrypt_email = var.admin_email
  letsencrypt_agree = true

  dns_challenge            = true
  dns_provider             = "cloudflare"
  dns_provider_credentials = "dns_cloudflare_api_token=${var.dns_cloudflare_api_token}"
  propagation_seconds      = 10
}

//This is just commented out for now, because ultimately I'd rather manage this from here than passing it in the other way
//data "nginxproxymanager_access_list" "access_list" {
//  id = var.internal_only ? module.nginx.internal_access_list_id : module.nginx.cloudflare_access_list_id
//}

resource "nginxproxymanager_proxy_host" "host" {
  domain_names = [var.domain_name]

  forward_scheme = var.forward_scheme
  forward_host   = var.service_ipv4
  forward_port   = var.service_port

  caching_enabled         = true
  allow_websocket_upgrade = true
  block_exploits          = true

  access_list_id = var.access_list_id

  certificate_id  = nginxproxymanager_certificate_letsencrypt.certificate.id
  ssl_forced      = false
  hsts_enabled    = false
  hsts_subdomains = false
  http2_support   = true
}