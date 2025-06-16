module "service_docker" {
  source = "../../modules/docker"
  
  container_name = "AA_Grafana"
  container_image = "grafana/grafana:12.0.1-security-01"
  attach_to_br0 = false
  attach_to_br1 = true
  br1_ipv4_addr = local.ip_address
  environment_vars = [
    "GF_AUTH_GENERIC_OAUTH_ENABLED=true",
    "GF_AUTH_GENERIC_OAUTH_NAME=authentik",
    "GF_AUTH_GENERIC_OAUTH_CLIENT_ID=${module.authentication.client_id}",
    "GF_AUTH_GENERIC_OAUTH_CLIENT_SECRET=${module.authentication.client_secret}",
    "GF_AUTH_GENERIC_OAUTH_SCOPES=openid profile email",
    "GF_AUTH_OAUTH_AUTO_LOGIN=true",
    "GF_AUTH_GENERIC_OAUTH_AUTH_URL=${local.auth_domain}/application/o/authorize/",
    "GF_AUTH_GENERIC_OAUTH_TOKEN_URL=${local.auth_domain}/application/o/token/",
    "GF_AUTH_GENERIC_OAUTH_API_URL=${local.auth_domain}/application/o/userinfo/",
    "GF_AUTH_SIGNOUT_REDIRECT_URL=${local.auth_domain}/application/o/${module.authentication.application_slug}/end-session/",
    "GF_SERVER_ROOT_URL=/",
    ]
}

/*
environment:
    # Optionally map user groups to Grafana roles
    GF_AUTH_GENERIC_OAUTH_ROLE_ATTRIBUTE_PATH: "contains(groups, 'Grafana Admins') && 'Admin' || contains(groups, 'Grafana Editors') && 'Editor' || 'Viewer'"
*/

module "service_dns" {
  source = "../../modules/dns"
  
  internal_only = true
  service_port = local.service_port
  zone_name = local.zone_name
  domain_name = local.domain_name
  access_list_id = var.access_list_id
  internal_host_ipv4 = local.nginx_ip
  service_ipv4 = local.ip_address
  admin_email = var.admin_email
  dns_cloudflare_api_token = var.cloudflare_token
  external_host_ipv4 = var.public_facing_ip
}

module "authentication" {
  source = "../../modules/oauth_auth"
  
  allowed_redirect_uris = [
    {
      matching_mode = "strict",
      url           = "https://${local.domain_name}/login/generic_oauth",
    }
  ]

  name = local.app_name
  
  create_access_group = true
  access_group_name = "terraform_${local.app_name}"
  user_to_add_to_access_group = var.admin_username
}