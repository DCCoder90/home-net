module "service_docker" {
  source = "../../modules/docker"
  
  container_name = "AAutomated_Test"
  container_image = "linuxserver/prowlarr:1.37.0"
  attach_to_br0 = false
  attach_to_br1 = true
  br1_ipv4_addr = local.ip_address
}

module "service_dns" {
  source = "../../modules/dns"
  
  internal_only = true
  service_port = 9000//local.service_port  - If using proxy_auth it needs to be the port to authentik
  zone_name = "mallett.family"
  domain_name = local.domain_name
  access_list_id = var.access_list_id
  internal_host_ipv4 = "192.168.4.2" //Port to Nginx
  service_ipv4 = local.authentik_ip
  admin_email = var.admin_email
  dns_cloudflare_api_token = var.cloudflare_token
  external_host_ipv4 = var.public_facing_ip
}

//Note: This will still require manually navigating to the service and setting up the credentials on the first run

module "authentication" {
  source = "../../modules/proxy_auth"
  
  internal_host = "http://${local.ip_address}:${local.service_port}"
  external_host = local.domain_name
  name = "prowlarr"
  username_attribute = "prowlarr_username"
  password_attribute = "prowlarr_password"
  create_access_group = true
  access_group_name = "terraform_prowlarr"
  user_to_add_to_access_group = var.admin_username
  access_group_attributes = jsonencode(
    {
      prowlarr_username: local.username,
      prowlarr_password: local.password
    }
  )
}