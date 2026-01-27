provider "docker" {
  host = "unix:///var/run/docker.sock"
}

provider "infisical" {
  auth = {
    universal = {
      client_id     = var.infiscal_client_id
      client_secret = var.infiscal_client_secret
    }
  }
}
