variable "infiscal_client_id" {
  type        = string
  sensitive   = true
  description = "Client ID for Infisical universal auth. Set as a sensitive variable in the HCP Terraform workspace."
}

variable "infiscal_client_secret" {
  type        = string
  sensitive   = true
  description = "Client secret for Infisical universal auth. Set as a sensitive variable in the HCP Terraform workspace."
}
