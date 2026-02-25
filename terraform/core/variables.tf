variable "nginx_proxy_address" { type = string }
variable "nginx_proxy_user" { type = string }
variable "nginx_proxy_pass" { type = string }

variable "tz" {
  type        = string
  description = "Timezone for services"
  default     = "America/Chicago" # Update as needed
}

variable "dns_admin_password" {
  type        = string
  sensitive   = true
  description = "Password for Technitium DNS"
}

variable "host_path_dns_config" {
  type    = string
  default = "/mnt/user/appdata/technitium"
}

variable "host_path_npm_data" {
  type    = string
  default = "/mnt/user/appdata/npm/data"
}

variable "host_path_npm_letsencrypt" {
  type    = string
  default = "/mnt/user/appdata/npm/letsencrypt"
}

# Variables provided by Terraform Cloud variable sets
variable "public_facing_ip" { type = string }
variable "vpn_user" { type = string }
variable "technitium_host" { type = string }
variable "network_admin_email" { type = string }
variable "network_admin_username" { type = string }

variable "vpn_pass" {
  type      = string
  sensitive = true
}

variable "technitium_password" {
  type      = string
  sensitive = true
}

variable "technitium_api_token" {
  type      = string
  sensitive = true
}

variable "infiscal_client_id" {
  type      = string
  sensitive = true
}

variable "infiscal_client_secret" {
  type      = string
  sensitive = true
}

variable "cloudflare_api_token" {
  type      = string
  sensitive = true
}
