output "client_id" {
  description = "Shared client ID for the OAuth2 provider"
  value       = authentik_provider_oauth2.name.client_id
}

output "client_secret" {
  description = "Shared client secret for the OAuth2 provider"
  value       = authentik_provider_oauth2.name.client_secret
}

output "application_slug" {
  description = "Slug of application"
  value       = authentik_application.name.slug
}

output "application_name" {
  description = "Name of application"
  value       = authentik_application.name.name
}

output "authorize_url" {
  value = data.authentik_provider_oauth2_config.oauth_config.authorize_url
}

output "issuer_url" {
  value = data.authentik_provider_oauth2_config.oauth_config.issuer_url
}

output "logout_url" {
  value = data.authentik_provider_oauth2_config.oauth_config.logout_url
}

output "token_url" {
  value = data.authentik_provider_oauth2_config.oauth_config.token_url
}

output "user_info_url" {
  value = data.authentik_provider_oauth2_config.oauth_config.user_info_url
}

output "jwks_url" {
  value = data.authentik_provider_oauth2_config.oauth_config.jwks_url
}

output "provider_info_url" {
  value = data.authentik_provider_oauth2_config.oauth_config.provider_info_url
}