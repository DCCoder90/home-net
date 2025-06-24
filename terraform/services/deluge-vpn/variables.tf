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

variable "service" {
  type = object({
    service_name = string
    image_name   = string
    icon         = optional(string)
    volumes      = optional(list(string))
    network = object({
      internal     = bool
      service_port = number
      networks = list(object({
        name       = string
        ip_address = string
      }))
    })
    env = list(string)
    mounts       = optional(list(string))
    capabilities = optional(object({
      add  = optional(list(string))
      drop = optional(list(string))
    }), {})
  })

  description = "Service configuration"
}