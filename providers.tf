provider "nginxproxymanager" {
  url      = var.nginx_proxy_address
  username = var.nginx_proxy_user
  password = var.nginx_proxy_pass
}

provider "portainer" {
}

//The Agent is running on the host that has docker, so we'll just connect directly to the socket
provider "docker" {
  host = "unix:///var/run/docker.sock"
}