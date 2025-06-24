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
  default     = null
}

variable "create_access_group" {
  type        = bool
  default     = false
  description = "If true, create new access group(s), otherwise use pre-created access group(s)"
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

variable "allowed_redirect_uris" {
  description = "List of redirect URIs for the OAuth2 provider"
  type        = list(map(string))
  default = [
    {
      matching_mode = "strict",
      url           = "http://localhost",
    }
  ]
}