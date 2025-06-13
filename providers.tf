provider "nginxproxymanager" {
  url = var.nginx_proxy_address
  username = var.nginx_proxy_user
  password = var.nginx_proxy_pass
}

provider "portainer" {
}