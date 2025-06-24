variable "internal_host" {
  description = "The IP and port used to connect the the application"
  type        = string
}

variable "external_host" {
  description = "The domain name used to access the application, (includes HTTPS)"
  type        = string
}

variable "name" {
  description = "Name of the application"
  type        = string
}

variable "description" {
  description = "Description of the application"
  type        = string
  default     = ""
}

variable "group" {
  description = "Group to assign the application to"
  type        = string
  default     = ""
}

variable "create_access_group" {
  type        = bool
  default     = false
  description = "If true, create new access group, otherwise use pre-created access group"
}

variable "access_group_name" {
  description = "Group to assign access to the application"
  type        = string
  default     = ""
}

variable "user_to_add_to_access_group" {
  type        = string
  description = "User to add to access group"
  default     = ""
}

variable "access_group_roles" {
  type        = list(string)
  description = "Roles to add to access group (if create_access_group = true)"
  default     = []
}

variable "access_group_attributes" {
  type        = string
  description = "JSON format expected. Use jsonencode() to pass objects. "
  default     = "{}"
}

variable "username_attribute" {
  description = "The attribute to use for the username in basic auth"
  type        = string
  default     = "username"
}

variable "password_attribute" {
  description = "The attribute to use for the password in basic auth"
  type        = string
  default     = "password"
}