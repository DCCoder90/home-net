data "infisical_projects" "home-net" {
  # This assumes a constant project slug. Consider making this configurable
  # via var.system if it needs to be more dynamic.
  slug = var.system.infisical.project
}

data "infisical_secrets" "secrets" {
  # Only fetch secrets if the stack configuration requests them.
  count = (var.service.secrets != null && length(var.service.secrets) > 0) || (try(var.service.auth.proxy.enabled, false)) ? 1 : 0
 
  env_slug     = var.system.infisical.environment
  workspace_id = data.infisical_projects.home-net.id
  # This path corresponds to where the root `secrets` module stores secrets.
  folder_path = var.system.infisical.folder
}

locals{
  # Create a list of environment variables from the secrets map.
  secret_envs = (var.service.secrets != null && length(var.service.secrets) > 0) ? [
    for key, value in var.service.secrets :
    format("%s=%s", key, data.infisical_secrets.secrets[0].secrets[value].value)
  ] : []
}