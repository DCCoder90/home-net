terraform {
  required_providers {
    infisical = {
      source = "infisical/infisical"
      version = "0.15.19"
    }
  }
}

data "infisical_projects" "home-net" {
  slug = var.project_slug
}

resource "random_password" "generated_secret" {
  for_each = toset(var.name)

  length           = 24
  special          = true
  override_special = "!#$%&*()-_=+[]{}<>:?"
}

resource "infisical_secret" "stored_secret" {
  for_each = local.generated_secrets_map

  name         = each.key
  value        = each.value
  env_slug     = "dev"
  workspace_id = data.infisical_projects.home-net.id
  folder_path  = "/generated/credentials"
}