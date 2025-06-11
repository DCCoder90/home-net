module "nginx" {
  source = "../nginx_config"
}

resource "nginxproxymanager_certificate_letsencrypt" "certificate" {
  domain_names = [var.domain_name]

  letsencrypt_email = var.admin_email
  letsencrypt_agree = true

  dns_challenge            = true
  dns_provider             = "cloudflare"
  dns_provider_credentials = var.dns_cloudflare_api_token
}

data "nginxproxymanager_access_list" "access_list" {
  id = var.internal_only ? module.nginx.outputs.internal_access_list_id : module.nginx.outputs.cloudflare_access_list_id
}

resource "nginxproxymanager_proxy_host" "host" {
  domain_names = [var.domain_name]

  forward_scheme = "https"
  forward_host   = var.domain_name
  forward_port   = 443

  caching_enabled         = true
  allow_websocket_upgrade = true
  block_exploits          = true

  access_list_id = data.nginxproxymanager_access_list.access_list.id

  locations = [
    {
      path           = "/"
      forward_scheme = "http"
      forward_host   = var.internal_host_ipv4 != "" ? var.internal_host_ipv4 : var.internal_host_ipv6
      forward_port   = var.service_port
    }
  ]

  certificate_id  = nginxproxymanager_certificate_letsencrypt.certificate.id
  ssl_forced      = true
  hsts_enabled    = false
  hsts_subdomains = false
  http2_support   = true
}