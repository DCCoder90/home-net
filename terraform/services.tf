module "arr_services" {
  source = "./modules/proxy_service_stack"

  admin_email      = var.network_admin_email
  admin_username   = var.network_admin_username
  cloudflare_token = var.cloudflare_api_token
  access_list_id   = module.nginx_conf.internal_access_list_id
  public_facing_ip = var.public_facing_ip
  stack            = local.stacks.arr_services
  zone_name        = local.system.zones.mallett
  proxy_ip         = local.system.proxy_ip
  authentik_ip     = local.system.authentik.ip_address
  authentik_port   = local.system.authentik.port
}

module "flaresolverr_service" {
  source   = "./modules/docker"

  container_name   = local.services.flaresolverr.service_name
  container_image  = local.services.flaresolverr.image_name
  attach_to_br0    = false
  attach_to_br1    = true
  br1_ipv4_addr    = local.services.flaresolverr.ip_address
  environment_vars = local.services.flaresolverr.env
}

module "delugevpn_service" {
  source   = "./services/deluge-vpn"
  service = local.services.deluge-vpn
  vpn_pass = var.vpn_pass
  vpn_user = var.vpn_user
}

module "karakeep_stack" {
  source = "./modules/docker-stack"
  stack  = local.stacks.karakeep
}
