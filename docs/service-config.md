## ⚙️ Standalone Service Definitions

Standalone service definition files define individual services or global configuration parameters that are directly consumed by specific Terraform modules. These are distinct from the more structured "Stack Configuration Files" and are typically used for services that don't fit into a multi-service stack pattern.

### File Location

All standalone service definition files must be placed in the `config/services/` directory. Each file should have a `.yaml` extension. Terraform will merge all top-level keys from all `.yaml` files found in this directory into a single `local.services` object.

### Structure

Each service definition file should define a top-level key that represents a unique identifier for the service or configuration. Under this key, you define the specific parameters for that service.

```yaml
your_service_identifier: # e.g., flaresolverr, deluge-vpn
  # Required: The canonical name of the service (used for Docker container name, etc.)
  service_name: "MyStandaloneService"
  # Required: The Docker image name and tag
  image_name: "myrepo/my-standalone-service:latest"
  # Required: The static IP address assigned to the container on its primary bridge network (e.g., br1)
  ip_address: "192.168.5.X"
  # Optional: The port the service listens on inside the container (used for web_ui, etc.)
  service_port: 8080
  # Optional: List of environment variables in "KEY=VALUE" format
  env:
    - "TZ=America/Chicago"
    - "SOME_VAR=some_value"
  # Optional: List of bind mounts in "host_path:container_path[:ro]" format
  mounts:
    - "/mnt/user/data:/data"
    - "/etc/localtime:/etc/localtime:ro"
  # Optional: Linux capabilities to add to the container
  capabilities:
    add:
      - "NET_ADMIN"
      - "SYS_PTRACE"
    drop:
      - "MKNOD"
      - "SETUID"
  # Optional: URL to an icon for the service (e.g., for Unraid UI)
  icon: "https://example.com/icon.png"
  # Optional: URL for the service's web UI (if not automatically derived)
  web_ui: "http://my-service.local:8080"
  # Other optional parameters that can be passed to the underlying 'docker' module:
  # container_restart: "unless-stopped" # Default: "unless-stopped"
  # container_user: "1000:100" # Default: null
  # container_ports: [] # List of objects {internal=number, external=optional(number), ...}
  # container_volumes: [] # List of objects for Docker-managed volumes
  # container_dns_servers: ["8.8.8.8"] # Default: ["8.8.8.8", "1.1.1.1"]
  # container_privileged_mode: false # Default: false
  # container_network_mode: "bridge" # Default: null

# Example of a non-Docker configuration:
authentik:
  admin-user: "your_authentik_admin_username"
```

### Field Meanings and Usage

*   **`your_service_identifier`**: This is the top-level key in your YAML file (e.g., `flaresolverr`, `deluge-vpn`). It acts as a unique identifier for this specific service configuration within the `local.services` map.
*   **`service_name`**: (Required) The human-readable name of the service, often used as the Docker container name.
*   **`image_name`**: (Required) The Docker image to pull, including its tag (e.g., `binhex/arch-delugevpn:2.2`).
*   **`ip_address`**: (Required) The static IPv4 address assigned to the container on its primary network interface (e.g., `br1`).
*   **`service_port`**: (Optional) The port the service listens on *inside* the container. This is used by modules that need to construct URLs or proxy configurations.
*   **`env`**: (Optional) A list of `KEY=VALUE` strings that will be set as environment variables inside the container.
*   **`mounts`**: (Optional) A list of bind mount strings in the format `host_path:container_path[:ro]`.
*   **`capabilities.add`**: (Optional) A list of Linux capabilities to add to the container (e.g., `NET_ADMIN`).
*   **`icon`**: (Optional) A URL to an icon for the service, often used by UI tools like Unraid.
*   **`web_ui`**: (Optional) A direct URL to the service's web interface. If not provided, some modules might attempt to construct one from `ip_address` and `service_port`.

### Special Cases

*   **`authentik`**: This top-level key in `config/services.yaml` is not a Docker service definition. Instead, it holds global configuration parameters for your Authentik instance (e.g., `admin-user`) that are consumed by other Terraform modules.