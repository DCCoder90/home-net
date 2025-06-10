resource "nginxproxymanager_proxy_host" "host" {
  domain_names = [var.External_Host]

  forward_scheme = "https"
  forward_host   = var.External_Host
  forward_port   = 443

  caching_enabled         = true
  allow_websocket_upgrade = true
  block_exploits          = true

  access_list_id = 1

  locations = [
    {
      path           = "/"
      forward_scheme = var.External_Scheme_Host
      forward_host   = var.External_Host
      forward_port   = var.Port
    }
  ]

  certificate_id  = 1
  ssl_forced      = true
  hsts_enabled    = false
  hsts_subdomains = false
  http2_support   = true
}

resource "nginxproxymanager_certificate_letsencrypt" "certificate" {
  domain_names = [var.External_Host]

  letsencrypt_email = var.Admin_Email
  letsencrypt_agree = true
  dns_challenge            = true
  dns_provider             = "cloudflare"
  dns_provider_credentials = var.dns_cloudflare_api_token
}