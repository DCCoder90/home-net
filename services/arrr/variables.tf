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

variable "services" {
  type = map(object({
    service_name   = string
    image_name     = string
    username       = string
    password       = string
    domain_name    = string
    service_port   = number
    ip_address     = string
    env            = list(string)
  }))
  description = "Map of services to be deployed with their configurations"
}