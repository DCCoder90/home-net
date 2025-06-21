variable "service" {
  type = object({
    icon      = optional(string)
    volumes  = optional(list(string))
    service_name = string
    image_name   = string
    service_port = number
    ip_address   = string
    env          = optional(list(string))
    mounts       = optional(list(string))
    capabilities = optional(object({
      add  = optional(list(string))
      drop = optional(list(string))
    }), {})
  })
  description = "Service and config"
}