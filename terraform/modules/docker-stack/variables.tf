variable "stack" {
  type = object({
    env     = optional(list(string))
    mounts  = optional(list(string))
    volumes = optional(list(string))
    generated_secrets = optional(list(string))
    zone_name = optional(string, null)
    networks = optional(map(object({
      internal = optional(bool, false)
      driver   = optional(string, "bridge")
      options  = optional(map(string), {})
    })))
    services = map(object({
      icon         = optional(string, "https://static-00.iconduck.com/assets.00/container-optimized-os-icon-512x512-9od9q084.png")
      volumes      = optional(list(string))
      dns          = object({
        enabled = optional(bool, false),
        internal = optional(bool, true),
        domain_name = optional(string)
      })
      auth = optional(object({
        enabled = optional(bool, false),
        proxy = optional(bool, false)
      }), {
        enabled = false,
        proxy = false
      })
      network = optional(object({
        internal = optional(bool, false)
        service_port = optional(number)
        ip_address   = optional(string)
        networks     = optional(list(string))
      }))
      service_name = string
      image_name   = string
      username     = optional(string, "")
      password     = optional(string, "")
      env          = optional(list(string))
      mounts       = optional(list(string))
      capabilities = optional(object({
        add  = optional(list(string))
        drop = optional(list(string))
      }), {})
    }))
  })
  description = "Map of services to be deployed with their configurations"
}

variable "system"{
  type = object({
      authentik              = object({
        ip_address             = string
        port                   = number
      })
      proxy_ip               = string
      admin_username         = string
      nginx_proxy_address    = string
      nginx_proxy_user       = string
      nginx_proxy_pass       = string
      cloudflare_api_token   = string
      network_admin_email    = string
      technitium_api_token   = string
      technitium_host        = string
      public_facing_ip       = string
      network_admin_username = string
      vpn_user               = string
      vpn_pass               = string
  })
}