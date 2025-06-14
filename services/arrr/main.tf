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
}

module "service_docker" {
  source = "../../modules/docker"
  
  container_name = "AAutomated_Test"
  container_image = "linuxserver/prowlarr:latest"
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
  service_ipv4 = local.ip_address
  admin_email = var.admin_email
  dns_cloudflare_api_token = var.cloudflare_token
  external_host_ipv4 = var.public_facing_ip
}