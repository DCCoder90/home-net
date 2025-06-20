module "arr_services" {
  source = "./modules/proxy_service_stack"

  admin_email      = var.network_admin_email
  admin_username   = var.network_admin_username
  cloudflare_token = var.cloudflare_api_token
  access_list_id   = module.nginx_conf.internal_access_list_id
  public_facing_ip = var.public_facing_ip
  services         = local.stacks.arr_services
  zone_name        = local.system.zones.mallett
  proxy_ip         = local.system.proxy_ip
  authentik_ip     = local.system.authentik.ip_address
  authentik_port   = local.system.authentik.port
}

/*

module "grafana_service" {
  source = "./services/grafana"

  admin_email      = var.network_admin_email
  admin_username   = var.network_admin_username
  cloudflare_token = var.cloudflare_api_token
  access_list_id   = module.nginx_conf.internal_access_list_id
  public_facing_ip = var.public_facing_ip
}
*/