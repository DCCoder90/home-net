module "arr_service" {
  source = "./services/arrr"

  admin_email      = var.network_admin_email
  cloudflare_token = var.cloudflare_api_token
  access_list_id   = module.nginx_conf.internal_access_list_id
  public_facing_ip = var.public_facing_ip
}
