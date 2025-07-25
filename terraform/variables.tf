variable "nginx_proxy_address" {
  type = string
}

variable "nginx_proxy_user" {
  type = string
}

variable "nginx_proxy_pass" {
  type = string
}

variable "cloudflare_api_token" {
  type        = string
  description = "Cloudflare API token"
}

variable "network_admin_email" {
  type        = string
  description = "Network admin email address"
}

variable "technitium_api_token" {
  type        = string
  description = "API token for technitium"
}

variable "technitium_host" {
  type        = string
  description = "Host for technitium"
}

variable "public_facing_ip" {
  type        = string
  description = "My public facing IP address"
}

variable "network_admin_username" {
  type        = string
  description = "Admin username"
}

variable "vpn_user" {
  type        = string
  description = "The username for the VPN connection."
  sensitive   = true
}

variable "vpn_pass" {
  type        = string
  description = "The password for the VPN connection."
  sensitive   = true
}

variable "infiscal_client_id" {
  type        = string
  description = "The password for the VPN connection."
  sensitive   = true
}

variable "infiscal_client_secret" {
  type        = string
  description = "The password for the VPN connection."
  sensitive   = true
}