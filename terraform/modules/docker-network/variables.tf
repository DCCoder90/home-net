variable "networks" {
  default = {}
  type = map(object({
    internal = optional(bool, false)
    driver   = optional(string, "bridge")
    options  = optional(map(string), {})
  }))
  description = "List of networks to create"
}