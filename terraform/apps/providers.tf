provider "nginxproxymanager" {
  url      = var.nginx_proxy_address
  username = var.nginx_proxy_user
  password = var.nginx_proxy_pass
}

provider "docker" {
  host = "unix:///var/run/docker.sock"
//  host     = "ssh://root@192.168.1.41:22"
//  ssh_opts = ["-o", "StrictHostKeyChecking=no", "-o", "UserKnownHostsFile=/dev/null"]
}

/*
provider "docker" {
}
*/

provider "technitium" {
  url   = var.technitium_host
  token = var.technitium_api_token
}

provider "cloudflare" {
  api_token = var.cloudflare_api_token
}

provider "authentik" {
}

provider "infisical" {
  auth = {
    universal = {
      client_id     = var.infiscal_client_id
      client_secret = var.infiscal_client_secret
    }
  }
}