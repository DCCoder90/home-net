variable "stack" {
  type = object({
    env               = optional(list(string))
    mounts            = optional(list(string))
    volumes           = optional(list(string))
    networks = optional(map(object({
      internal = optional(bool, false)
      driver   = optional(string, "bridge")
      options  = optional(map(string), {})
    })))
    services = map(object({
      service_name = string
      image_name   = string
      icon         = optional(string, "https://static-00.iconduck.com/assets.00/container-optimized-os-icon-512x512-9od9q084.png")
      network_mode = optional(string)
      description  = optional(string, "")
      commands     = optional(list(string))
      mounts       = optional(list(string))
      enable_gpu   = optional(bool, false)
      env          = optional(list(string))
      volumes      = optional(list(string)),
      secrets      = optional(map(string)),  
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
      enabled = optional(bool, false)
      group   = optional(string, "Uncategorized")
      proxy = optional(object({
        enabled     = optional(bool, false)
        user_secret = optional(string)
        pass_secret = optional(string)
      }), {})
      oauth = optional(object({
        enabled       = optional(bool, false),
        keys          = optional(map(string), {}),
        scopes        = optional(list(string)),
        redirect_uris = optional(list(string))
      }), {})
    }), {})
      network = optional(object({
        internal     = optional(bool, false)
        service_port = optional(number)
        networks = optional(list(object({
          name       = string
          ip_address = optional(string)
        })), [])
      }))
    }))
  })

  description = "Map of services to be deployed with their configurations"
}

variable "system" {
  type        = any
  description = "A dynamic object containing system-wide configuration variables. Its attributes are accessed via lookup() for flexibility."
}