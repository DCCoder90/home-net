module "stacks" {
  for_each = local.stacks
  source = "./modules/docker-stack"

  stack  = each.value
  system = local.system
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