# Core Infrastructure - Technitium DNS
# Managed explicitly outside of the standard service loop for stability.

resource "docker_container" "technitium" {
  name  = "technitium_dns"
  image = "technitium/dns-server:latest"

  # Hostname from your config
  hostname = "technitium-dns"

  # Vital for core infra
  restart = "unless-stopped"

  # Environment Variables
  env = [
    "TZ=America/Chicago", # Update if your ${TZ} is different
    "DNS_SERVER_ADMIN_PASSWORD=${var.technitium_password}"
  ]

  # Volumes
  # Note: Your Portainer config only listed /etc/dns. 
  # If you have a separate volume for /var/lib/dns, add it here too.
  volumes {
    host_path      = "/mnt/user/DataStore/technitium/config" # Replace with your actual ${DNS_CONFIG_PATH}
    container_path = "/etc/dns"
  }

  # Network Configuration (Static IP on br1)
  networks_advanced {
    name         = "br1"
    ipv4_address = "192.168.4.X" # REPLACE THIS with your actual ${DNS_IP}
  }

  # Lifecycle ensures we don't destroy the container if minor attributes change
  lifecycle {
    ignore_changes = [
      # If you manage DNS zones via the UI, you might want to ignore volume changes or specific envs
    ]
  }
}
