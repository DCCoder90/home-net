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
