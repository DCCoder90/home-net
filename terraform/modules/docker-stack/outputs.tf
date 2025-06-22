output "generated_secrets" {
  description = "A map of the dynamically generated secrets for the stack."
  value       = local.generated_secrets
  sensitive   = true
}