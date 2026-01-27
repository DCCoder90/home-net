output "proxy_ip" {
  value       = local.system.proxy_ip
  description = "The internal IP of the Nginx Proxy, defined in system.yaml"
}

output "internal_access_list_id" {
  value       = module.nginx_conf.internal_access_list_id
  description = "ID for Internal Only Access List"
}

output "cloudflare_access_list_id" {
  value       = module.nginx_conf.cloudflare_access_list_id
  description = "ID for Cloudflare Access List"
}
