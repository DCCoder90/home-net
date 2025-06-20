variable "admin_email" {
  description = "Admin email address"
}

variable "admin_username" {
  description = "Admin username"
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