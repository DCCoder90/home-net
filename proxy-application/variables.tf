variable "Internal_Host" {
    description = "Internal host for Application"
    type        = string
}

variable "External_Host" {
    description = "External host for Application"
}

variable "External_Scheme_Host"{
    description = "value for the external scheme host"
    type        = string
    default     = "http"
}

variable "Application_Name" {
    description = "Name of the Application"
    type        = string
}

variable "Port"{
    description = "Port for the Application"
    type        = number
    default     = 80
}

variable "Application_Group" {
    description = "Group to assign the application to"
    type        = string
    default     = ""
}

variable "Application_Description"{
    description = "Description of the Application"
    type        = string
    default     = ""
}

variable "Username_Attribute" {
    description = "Username attribute for Basic Auth"
    type        = string
    default     = ""
}

variable "Password_Attribute" {
    description = "Password attribute for Basic Auth"
    type        = string
    default     = ""
}

variable "Admin_Email" {
    description = "Email address for the admin user"
    type        = string
    default     = ""
}