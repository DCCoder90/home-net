data "infisical_projects" "home-net" {
  slug = local.system.infisical.project
}

data "infisical_secrets" "secrets" {
  env_slug     = local.system.infisical.environment
  workspace_id = data.infisical_projects.home-net.id
  folder_path  = local.system.infisical.folder
}
