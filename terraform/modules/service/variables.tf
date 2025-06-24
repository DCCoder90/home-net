variable "service" {
  type = object({
    icon         = optional(string, "https://static-00.iconduck.com/assets.00/container-optimized-os-icon-512x512-9od9q084.png")
    volumes      = optional(list(string))
    service_name = string
    network_mode = optional(string)
    description  = optional(string, "")
    image_name   = string
    username     = optional(string, "")
    password     = optional(string, "")
    commands     = optional(list(string))
    capabilities = optional(object({
      add  = optional(list(string))
      drop = optional(list(string))
    }), {})
    dns = object({
      enabled     = optional(bool, false),
      internal    = optional(bool, true),
      domain_name = optional(string)
    })
    auth = optional(object({
      enabled = optional(bool, false),
      proxy   = optional(bool, false),
      group   = optional(string, "Uncategorized"),
      oauth = optional(object({
        enabled       = optional(bool, false),
        keys          = optional(map(string), {}),
        scopes        = optional(list(string)),
        redirect_uris = optional(list(string))
      }), {})
    }))
    network = optional(object({
      internal     = optional(bool, false)
      service_port = optional(number)
      networks = optional(list(object({
        name         = string
        ipv4_address = optional(string)
      })), [])
    }))
  })

  description = "Map of services to be deployed with their configurations"
}

variable "system" {
  type        = any
  description = "A dynamic object containing system-wide configuration variables. Its attributes are accessed via lookup() for flexibility."
}