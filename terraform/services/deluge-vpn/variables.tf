variable "vpn_user"{
    type = string
    description = "The username for the VPN connection."
    sensitive = true
}

variable "vpn_pass"{
    type = string
    description = "The password for the VPN connection."
    sensitive = true
}

variable "service" {
  type = object({
    icon      = optional(string)
    volumes  = optional(list(string))
    service_name = string
    image_name   = string
    username     = optional(string, "")
    password     = optional(string, "")
    service_port = number
    ip_address   = string
    env          = optional(list(string))
    mounts       = optional(list(string))
  })
  description = "Service and config"
}