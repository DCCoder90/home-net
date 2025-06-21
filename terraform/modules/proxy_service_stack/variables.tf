variable "cloudflare_token" {
  type        = string
  description = "Cloudflare API token"
}

variable "admin_email" {
  type        = string
  description = "Network admin email address"
}

variable "admin_username" {
  type        = string
  description = "Network admin username"
}

variable "access_list_id" {
  type        = string
  description = ""
}

variable "public_facing_ip" {
  type        = string
  description = ""
}

variable "stack" {
  type = object({
    env      = optional(list(string))
    mounts   = optional(list(string))
    volumes  = optional(list(string))
    services = map(object({
      internal     = optional(bool, false)
      icon         = optional(string)
      service_name = string
      image_name   = string
      domain_name  = string
      service_port = number
      ip_address   = string
      env          = optional(list(string))
      mounts       = optional(list(string))
    }))
  })
  description = "Map of services to be deployed with their configurations"
}

variable "zone_name" {
  type        = string
  description = "DNS zone name for the services"
}

variable "proxy_ip" {
  type        = string
  description = "IP address of the proxy server"
}

variable "authentik_ip" {
  type        = string
  description = "IP address of the Authentik service"
}

variable "authentik_port" {
  type        = number
  description = "Port on which the Authentik service is running"
}