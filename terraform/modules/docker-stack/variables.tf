variable "stack" {
  type = object({
    env     = optional(list(string))
    mounts  = optional(list(string))
    volumes = optional(list(string))
    generated_secrets = optional(list(string))
    networks = optional(map(object({
      internal = optional(bool, false)
      driver   = optional(string, "bridge")
      options  = optional(map(string), {})
    })))
    services = map(object({
      icon         = optional(string, "https://static-00.iconduck.com/assets.00/container-optimized-os-icon-512x512-9od9q084.png")
      volumes      = optional(list(string))
      service_name = string
      image_name   = string
      username     = optional(string, "")
      password     = optional(string, "")
      service_port = optional(number, null)
      ip_address   = optional(string, "")
      env          = optional(list(string))
      mounts       = optional(list(string))
      capabilities = optional(object({
        add  = optional(list(string))
        drop = optional(list(string))
      }), {})
      networks = optional(list(string))
    }))
  })
  description = "Map of services to be deployed with their configurations"
}