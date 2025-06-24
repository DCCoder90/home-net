output "client_id" {
  description = "The OAuth client ID for the service, if OAuth is enabled."
  value       = { for k, m in module.oauth_authentication : k => m.client_id }
  sensitive   = true
}

output "client_secret" {
  description = "The OAuth client secret for the service, if OAuth is enabled."
  value       = { for k, m in module.oauth_authentication : k => m.client_secret }
  sensitive   = true
}

output "well_known_url" {
  description = "The well-known URL for the OAuth service, if OAuth is enabled."
  value       = { for k, m in module.oauth_authentication : k => m.provider_info_url }
  sensitive   = true
}