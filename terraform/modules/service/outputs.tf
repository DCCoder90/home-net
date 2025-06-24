output "generated_secrets" {
  description = "A map of the dynamically generated secrets for the stack."
  value       = local.generated_secrets
  sensitive   = true
}

output "client_id" {
  description = "The OAuth client ID for the service, if OAuth is enabled."
  value       = module.oauth_authentication.client_id
  sensitive   = true
}

output "client_secret" {
  description = "The OAuth client secret for the service, if OAuth is enabled."
  value       = module.oauth_authentication.client_secret
  sensitive   = true
}

output "well_known_url" {
  description = "The well-known URL for the OAuth service, if OAuth is enabled."
  value       = module.oauth_authentication.provider_info_url
  sensitive   = true
}