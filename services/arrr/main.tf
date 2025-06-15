variable "cloudflare_token" {
  type = string
  description = "Cloudflare API token"
}

variable "admin_email" {
  type = string
  description = "Network admin email address"
}

variable "access_list_id" {
  type = string
  description = ""
}

variable "public_facing_ip" {
  type = string
  description = ""
}

locals {
  domain_name="test.mallett.family"
  ip_address = "192.168.5.13"
  username = ""
  password = ""
  authentik_ip = ""
}

module "service_docker" {
  source = "../../modules/docker"
  
  container_name = "AAutomated_Test"
  container_image = "linuxserver/prowlarr:2.0.0-nightly"
  attach_to_br0 = false
  attach_to_br1 = true
  br1_ipv4_addr = local.ip_address
}

module "service_dns" {
  source = "../../modules/dns"
  
  internal_only = true
  service_port = 9696
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
  
  internal_host = local.ip_address
  external_host = local.domain_name
  name = "prowlarr"
  username_attribute = "prowlarr_username"
  password_attribute = "prowlarr_password"
  create_access_group = true
  access_group_name = "terraform_prowlarr"
  user_to_add_to_access_group = "MyAwesomeUser!"
  access_group_attributes = jsonencode(
    {
      prowlarr_username: local.username,
      prowlarr_password: local.password
    }
  )
}