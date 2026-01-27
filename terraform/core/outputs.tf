output "core_service_ips" {
  description = "IP addresses of core infrastructure services, for reference."
  value = {
    technitium = try(local.core_configs.technitium.networks[0].ip_address, "unknown")
    npm        = try(local.core_configs.nginx_proxy_manager.networks[0].ip_address, "unknown")
    authentik  = try(local.core_configs.authentik.services.server.networks[0].ip_address, "unknown")
  }
}
