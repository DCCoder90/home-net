variable "dccoder_email" {
  description = "My email address"
  default     = ""
}

variable "users" {
  description = "List of users to create in authentik"
  type = list(object({
    username = string
    email    = string
    password = string
  }))
  default = []

}