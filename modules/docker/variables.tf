variable "container_name"{
    type = string
    description = "Name for the container to be created"
}

variable "container_image"{
    type = string
    description = "Name and tag of the image to use (ex. ubuntu:latest)"
}

variable "environment_vars"{
    type=set(string)
    description = "Environment variables to set in the form of KEY=VALUE"
    default = null
}

variable "attach_to_br1"{
    type = bool
    description = "Attach to br1 network?"
    default = false
}

variable "attach_to_br0"{
    type = bool
    description = "Attach to br0 network?"
    default = false
}

variable "br1_ipv4_addr"{
    type = string
    description = "IPv4 address to assign on network br1"
    default = null
}

variable "br0_ipv4_addr"{
    type = string
    description = "IPv4 address to assign on network br0"
    default = null
}

variable "container_restart"{
    type = string
    description = "The restart policy for the container. Must be one of 'no', 'on-failure', 'always', 'unless-stopped'."
    default = "unless-stopped"
}

variable "container_user" {
    type = string
    default = null
    description = "User used for run the first process. Format is user or user:group which user and group can be passed literraly or by name."
}

variable "container_ports" {
  description = "List of port mappings for the container. Each object can have 'internal' (required, number), 'external' (optional, number), 'ip' (optional, string), and 'protocol' (optional, string, defaults to 'tcp')."
  type = list(object({
    internal = number
    external = optional(number)
    ip       = optional(string)
    protocol = optional(string)
  }))
  default = []
}

variable "container_volumes" {
  description = "List of volume mounts. Each object defines 'container_path' and EITHER 'host_path' (for bind mounts) OR 'volume_name' (for Docker-managed volumes). The 'manage_volume_lifecycle' flag determines if named volumes are created by this module."
  type = list(object({
    container_path = string
    read_only      = optional(bool, false)

    # For bind mounts
    host_path      = optional(string)

    # For Docker-managed named volumes
    volume_name    = optional(string)
    manage_volume_lifecycle = optional(bool, true) # If true (default), module creates/manages the named volume. If false, assumes volume_name exists.
    # Options for creating the named volume (used only if volume_name is set AND manage_volume_lifecycle is true)
    driver         = optional(string)
    driver_opts    = optional(map(string))
    labels         = optional(map(string))
  }))
  default     = []

  validation {
    condition = alltrue([
      for vol in var.container_volumes : (vol.host_path == null || vol.volume_name == null)
    ])
    error_message = "A volume entry cannot specify both 'host_path' and 'volume_name'."
  }
  validation {
    condition = alltrue([
      for vol in var.container_volumes : (vol.host_path != null || vol.volume_name != null)
    ])
    error_message = "A volume entry must specify either 'host_path' or 'volume_name'."
  }
}

variable "container_dns_servers" {
  type        = list(string)
  description = "List of DNS servers to use for the container."
  default     = ["8.8.8.8", "1.1.1.1"]
}

variable "container_privileged_mode" {
  type        = bool
  description = "Run the container in privileged mode. Grants the container almost all the capabilities of the host machine."
  default     = false
}