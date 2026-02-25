# --- Technitium DNS ---
resource "docker_image" "technitium" {
  name         = "technitium/dns-server:latest"
  keep_locally = true
}

resource "docker_container" "technitium_dns" {
  name     = "technitium_dns"
  image    = docker_image.technitium.image_id
  restart  = "unless-stopped"
  hostname = "technitium-dns"

  env = [
    "TZ=${var.tz}",
    "DNS_SERVER_ADMIN_PASSWORD=${var.dns_admin_password}"
  ]

  volumes {
    host_path      = var.host_path_dns_config
    container_path = "/etc/dns"
  }

  networks_advanced {
    name         = data.docker_network.br1.name
    ipv4_address = local.system.dns_ip
  }

  lifecycle {
    ignore_changes = [
      image,
      labels
    ]
  }
}

# --- Nginx Proxy Manager ---
resource "docker_image" "nginx_proxy" {
  name         = "jc21/nginx-proxy-manager:latest"
  keep_locally = true
}

resource "docker_container" "nginx_proxy" {
  name    = "jc21/nginx-proxy-manager"
  image   = docker_image.nginx_proxy.image_id
  restart = "unless-stopped"

  # Standard Ports
  ports {
    internal = 80
    external = 80
  }
  ports {
    internal = 81
    external = 81
  }
  ports {
    internal = 443
    external = 443
  }

  volumes {
    host_path      = var.host_path_npm_data
    container_path = "/data"
  }

  volumes {
    host_path      = var.host_path_npm_letsencrypt
    container_path = "/etc/letsencrypt"
  }

  networks_advanced {
    name         = data.docker_network.br1.name
    ipv4_address = local.system.proxy_ip
  }

  lifecycle {
    ignore_changes = [image, labels]
  }
}
