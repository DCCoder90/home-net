variable "cloudflare_token" {
  type = string
  description = "Cloudflare API token"
}

variable "admin_email" {
  type = string
  description = "Network admin email address"
}

variable "admin_username" {
  type = string
  description = "Network admin username"
}

variable "access_list_id" {
  type = string
  description = ""
}

variable "public_facing_ip" {
  type = string
  description = ""
}