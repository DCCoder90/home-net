module "delugevpn_service" {
  source = "../../modules/docker"

  icon            = var.service.icon
  web_ui          = "http://${var.service.network.networks[0].ip_address}:${var.service.network.service_port}"
  container_name  = var.service.service_name
  container_image = var.service.image_name
  networks        = var.service.network.networks
  environment_vars = concat(var.service.env, [
    "VPN_USER=${var.vpn_user}",
    "VPN_PASS=${var.vpn_pass}"
  ])
  mounts                 = var.service.mounts
  container_capabilities = var.service.capabilities
}