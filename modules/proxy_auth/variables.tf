variable "internal_host" {
  description = "If true, the application is only accessible internally"
  type        = bool
}

variable "external_host" {
  description = "If true, the application is only accessible internally"
  type        = bool
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

variable "access_group" {
  description = "Group to assign access to the application"
  type        = string
  default     = ""
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