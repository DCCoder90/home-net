/*
module "arr_service" {
  source = "./services/arrr"

  admin_email      = var.network_admin_email
  admin_username   = var.network_admin_username
  cloudflare_token = var.cloudflare_api_token
  access_list_id   = module.nginx_conf.internal_access_list_id
  public_facing_ip = var.public_facing_ip
}
*/

module "grafana_service" {
  source = "./services/grafana"

  admin_email      = var.network_admin_email
  admin_username   = var.network_admin_username
  cloudflare_token = var.cloudflare_api_token
  access_list_id   = module.nginx_conf.internal_access_list_id
  public_facing_ip = var.public_facing_ip
}