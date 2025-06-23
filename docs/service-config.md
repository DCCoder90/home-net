## ⚙️ Standalone Service Definitions

Standalone service definition files provide a flexible way to define individual services or global configuration parameters that are consumed directly by modules in the root `services.tf` file. This is distinct from the more structured "Stack Configuration Files" and is used for services that have unique requirements or don't fit into a multi-service stack.

### File Location

All standalone service definition files must be placed in the `config/services/` directory. Each file should have a `.yaml` or `.yml` extension. Terraform will merge all top-level keys from all files found in this directory into a single `local.services` object.


### Structure

The structure of each top-level key in these files is **not standardized**. It is determined entirely by the specific Terraform module that consumes it. 

Below are examples of different structures found in this project.

### Example 1: A Standalone Docker Service (`deluge-vpn`)

This service is consumed by the `delugevpn_service` module in `services.tf`, which is a wrapper around the generic `docker` module.

```yaml
deluge-vpn:
  service_name: "deluge-vpn"
  image_name: "binhex/arch-delugevpn:2.2"
  ip_address: "192.168.5.20"
  env:
    - "VPN_ENABLED=yes"
    - "VPN_PROV=custom"
    # ... other env vars
  mounts:
    - "/etc/localtime:/etc/localtime:ro"
  capabilities:
    add:
      - "NET_ADMIN"
```

### Example 2: A Global Configuration Object (`authentik`)

This key does not define a deployable service. Instead, it holds configuration parameters for the Authentik provider and other modules that need to interact with your Authentik instance.

```yaml
authentik:
  admin-user: "akadmin"
```

This structure provides a flexible way to define individual services or global parameters that are directly integrated into your Terraform deployment.

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