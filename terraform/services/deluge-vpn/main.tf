module "delugevpn_service" {
  source   = "../../modules/docker"

  icon = var.service.icon
  web_ui = "http://${var.service.ip_address}:${var.service.service_port}"
  container_name   = var.service.service_name
  container_image  = var.service.image_name
  attach_to_br0    = false
  attach_to_br1    = true
  br1_ipv4_addr    = var.service.ip_address
    environment_vars = concat(var.service.env, [
    "VPN_USER=${var.vpn_user}",
    "VPN_PASS=${var.vpn_pass}"
  ])
  mounts           = var.service.mounts
}