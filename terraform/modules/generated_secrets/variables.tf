variable "name"{
    type = list(string)
    description = "List of secret names to generate"
}

variable "project_slug" {
    type = string  
    description = "Slug of project to store secret in"
    default = "home-net-ln-sy"
  
}