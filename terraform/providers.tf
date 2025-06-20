provider "nginxproxymanager" {
  url      = var.nginx_proxy_address
  username = var.nginx_proxy_user
  password = var.nginx_proxy_pass
}

provider "docker" {
  host = "unix:///var/run/docker.sock"
}

provider "technitium" {
  url   = var.technitium_host
  token = var.technitium_api_token
}

provider "cloudflare" {
  api_token = var.cloudflare_api_token
}