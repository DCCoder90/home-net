output "generated_secrets_list" {
  description = "A list of objects, each containing the name and generated value of a secret."
  value       = local.generated_secret_list
  sensitive   = true
}

output "generated_secrets_map" {
  description = "A map of generated secret names to their values."
  value       = local.generated_secrets_map
  sensitive   = true
}