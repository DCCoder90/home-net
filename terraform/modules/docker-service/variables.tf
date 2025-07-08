variable "service" {
  type = object({
    # --- Core Service Definition ---
    service_name = string
    image_name   = string
    description  = optional(string, "")
    icon         = optional(string, "https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png/docker.png")

    # --- Execution & Runtime ---
    commands     = optional(list(string))
    network_mode = optional(string)
    enable_gpu   = optional(bool, false)
    capabilities = optional(object({
      add  = optional(list(string))
      drop = optional(list(string))
    }), {})

    # --- Storage ---
    mounts  = optional(list(string))
    volumes = optional(list(string))

    # --- Environment & Secrets ---
    env               = optional(list(string))
    secrets = optional(map(string))

    # --- Networking & DNS ---
    network = optional(object({ 
      internal     = optional(bool, false)
      service_port = optional(number)
      networks = optional(list(object({
        name       = string
        ip_address = optional(string)
      })), [])
    }))
    dns = object({
      enabled     = optional(bool, false),
      internal    = optional(bool, true),
      domain_name = optional(string)
    })

    # --- Authentication ---
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
  })

  description = "Service and it's configuration"
}

variable "system" {
  type        = any
  description = "A dynamic object containing system-wide configuration variables. Its attributes are accessed via lookup() for flexibility."
}