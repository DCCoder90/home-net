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

  labels {
    label = "net.unraid.docker.icon"
    value = "https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png/technitium.png"
  }

  labels {
    label = "net.unraid.docker.webui"
    value = "http://[IP]:[PORT:5380]"
  }

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
      labels,
      log_opts
    ]
  }
}

# --- Nginx Proxy Manager ---
resource "docker_image" "nginx_proxy" {
  name         = "jc21/nginx-proxy-manager:latest"
  keep_locally = true
}

resource "docker_container" "nginx_proxy" {
  name    = "Nginx-Proxy-Manager-Official"
  image   = docker_image.nginx_proxy.image_id
  restart = "unless-stopped"

  env = [
    "TZ=${var.tz}",
  ]

  labels {
    label = "net.unraid.docker.icon"
    value = "https://cdn.jsdelivr.net/gh/homarr-labs/dashboard-icons/png/nginx-proxy-manager.png"
  }

  labels {
    label = "net.unraid.docker.webui"
    value = "http://[IP]:[PORT:81]"
  }

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

  volumes {
    host_path      = var.host_path_npm_logs
    container_path = "/var/log"
  }

  networks_advanced {
    name         = data.docker_network.br1.name
    ipv4_address = local.system.proxy_ip
  }

  lifecycle {
    ignore_changes = [image, labels, env, log_opts, ports, volumes]
  }
}
