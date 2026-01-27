# -------------------------------------------------------------------
# Technitium DNS Server
# -------------------------------------------------------------------
module "technitium" {
  source = "git@github.com:DCCoder90/home-tf-modules.git//docker?ref=1.0.0"

  container_name  = local.core_configs.technitium.service_name
  container_image = local.core_configs.technitium.image_name
  environment_vars = toset(concat(
    coalesce(local.core_configs.technitium.env, []),
    [for k, v in coalesce(local.core_configs.technitium.secrets, {}) :
      "${k}=${data.infisical_secrets.secrets.secrets[v].value}"
    ]
  ))
  mounts   = local.core_configs.technitium.mounts
  networks = local.core_configs.technitium.networks
}

# -------------------------------------------------------------------
# Nginx Proxy Manager
# -------------------------------------------------------------------
module "nginx_proxy_manager" {
  source = "git@github.com:DCCoder90/home-tf-modules.git//docker?ref=1.0.0"

  container_name  = local.core_configs.nginx_proxy_manager.service_name
  container_image = local.core_configs.nginx_proxy_manager.image_name
  environment_vars = toset(coalesce(local.core_configs.nginx_proxy_manager.env, []))
  mounts   = local.core_configs.nginx_proxy_manager.mounts
  networks = local.core_configs.nginx_proxy_manager.networks
}

# -------------------------------------------------------------------
# Authentik PostgreSQL Database
# -------------------------------------------------------------------
module "authentik_postgres" {
  source = "git@github.com:DCCoder90/home-tf-modules.git//docker?ref=1.0.0"

  container_name  = local.core_configs.authentik.services.postgres.service_name
  container_image = local.core_configs.authentik.services.postgres.image_name
  environment_vars = toset(concat(
    coalesce(local.core_configs.authentik.services.postgres.env, []),
    [for k, v in coalesce(local.core_configs.authentik.services.postgres.secrets, {}) :
      "${k}=${data.infisical_secrets.secrets.secrets[v].value}"
    ]
  ))
  mounts   = coalesce(local.core_configs.authentik.services.postgres.mounts, [])
  networks = local.core_configs.authentik.services.postgres.networks
}

# -------------------------------------------------------------------
# Authentik Redis
# -------------------------------------------------------------------
module "authentik_redis" {
  source = "git@github.com:DCCoder90/home-tf-modules.git//docker?ref=1.0.0"

  container_name  = local.core_configs.authentik.services.redis.service_name
  container_image = local.core_configs.authentik.services.redis.image_name
  commands        = local.core_configs.authentik.services.redis.commands
  environment_vars = toset(coalesce(local.core_configs.authentik.services.redis.env, []))
  mounts   = coalesce(local.core_configs.authentik.services.redis.mounts, [])
  networks = local.core_configs.authentik.services.redis.networks
}

# -------------------------------------------------------------------
# Authentik Server
# -------------------------------------------------------------------
module "authentik_server" {
  source = "git@github.com:DCCoder90/home-tf-modules.git//docker?ref=1.0.0"

  depends_on = [module.authentik_postgres, module.authentik_redis]

  container_name  = local.core_configs.authentik.services.server.service_name
  container_image = local.core_configs.authentik.services.server.image_name
  container_privileged_mode = true
  container_user  = "1000"
  environment_vars = toset(concat(
    coalesce(local.core_configs.authentik.services.server.env, []),
    [for k, v in coalesce(local.core_configs.authentik.services.server.secrets, {}) :
      "${k}=${data.infisical_secrets.secrets.secrets[v].value}"
    ]
  ))
  mounts   = coalesce(local.core_configs.authentik.services.server.mounts, [])
  networks = local.core_configs.authentik.services.server.networks
}

# -------------------------------------------------------------------
# Authentik Worker
# -------------------------------------------------------------------
module "authentik_worker" {
  source = "git@github.com:DCCoder90/home-tf-modules.git//docker?ref=1.0.0"

  depends_on = [module.authentik_postgres, module.authentik_redis]

  container_name  = local.core_configs.authentik.services.worker.service_name
  container_image = local.core_configs.authentik.services.worker.image_name
  container_user  = "root"
  commands        = local.core_configs.authentik.services.worker.commands
  environment_vars = toset(concat(
    coalesce(local.core_configs.authentik.services.worker.env, []),
    [for k, v in coalesce(local.core_configs.authentik.services.worker.secrets, {}) :
      "${k}=${data.infisical_secrets.secrets.secrets[v].value}"
    ]
  ))
  mounts   = coalesce(local.core_configs.authentik.services.worker.mounts, [])
  networks = coalesce(local.core_configs.authentik.services.worker.networks, [])
}
