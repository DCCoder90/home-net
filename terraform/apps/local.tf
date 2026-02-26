locals {
  # Load system config
  system = merge(
    yamldecode(file("${path.module}/../../config/system.yaml")),
    {
      # We override the proxy_ip from system.yaml with the authoritative one from Core state
      # This ensures if you change the IP in Core, Apps knows about it automatically.
      proxy_ip = data.terraform_remote_state.core.outputs.proxy_ip

      # Pass through variables
      nginx_proxy_address    = var.nginx_proxy_address
      nginx_proxy_user       = var.nginx_proxy_user
      nginx_proxy_pass       = var.nginx_proxy_pass
      cloudflare_api_token   = var.cloudflare_api_token
      network_admin_email    = var.network_admin_email
      technitium_api_token   = var.technitium_api_token
      technitium_host        = var.technitium_host
      public_facing_ip       = var.public_facing_ip
      network_admin_username = var.network_admin_username
      vpn_user               = var.vpn_user
      vpn_pass               = var.vpn_pass
    }
  )

  # Load Stacks
  stacks = merge([
    for f in fileset("${path.module}/../../config/stacks", "*.yaml") :
    yamldecode(file("${path.module}/../../config/stacks/${f}"))
  ]...)

  # Load Services
  services = merge([
    for f in fileset("${path.module}/../../config/services", "*.yaml") :
    yamldecode(file("${path.module}/../../config/services/${f}"))
  ]...)
}
