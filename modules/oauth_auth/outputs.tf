output "client_id" {
  description = "Shared client ID for the OAuth2 provider"
  value       = authentik_provider_oauth2.name.client_id
}

output "client_secret" {
  description = "Shared client secret for the OAuth2 provider"
  value       = authentik_provider_oauth2.name.client_secret
}