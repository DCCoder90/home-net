## ⚙️ Standalone Service Definitions

Standalone service definition files provide a flexible way to define individual services that are consumed directly by modules in the root `services.tf` file. This is distinct from the more structured "Stack Configuration Files" and is used for services that have unique requirements or don't fit into a multi-service stack.

### File Location

All standalone service definition files must be placed in the `config/services/` directory. Each file should have a `.yaml` extension. Pulumi will merge all top-level keys from all files found in this directory into a single services map.

### Structure

Each standalone service file defines a single top-level key that acts as a unique identifier. The structure mirrors the service-level schema used inside stacks, with `dns`, `network`, `auth`, `secrets`, and other fields at the top level.

### Example 1: A Standalone Docker Service (`flaresolverr`)

```yaml
flaresolverr:
  dns:
    enabled: false
  network:
    internal: true
    service_port: 8112
    networks:
      - name: "br1"
        ip_address: "192.168.5.29"
  auth:
    enabled: false
  service_name: "flaresolverr"
  image_name: "flaresolverr/flaresolverr:v3.4.6"
```

### Example 2: A Standalone Docker Service with Secrets and Capabilities (`deluge-vpn`)

```yaml
deluge-vpn:
  dns:
    enabled: false
  network:
    internal: true
    service_port: 8112
    networks:
      - name: "br1"
        ip_address: "192.168.5.28"
  auth:
    enabled: false
    group: "Arr"
  secrets:
    VPN_USER: "VPN_USER"
    VPN_PASS: "VPN_PASS"
  service_name: "deluge-vpn"
  image_name: "binhex/arch-delugevpn:2.2"
  mounts:
    - "/etc/localtime:/etc/localtime:ro"
    - "/mnt/user/Arr/deluge-data:/config"
    - "/mnt/user/Downloads:/data/downloads"
    - "/mnt/user/Arr/deluge-data/openvpn:/config/openvpn"
  capabilities:
    add:
      - "CAP_NET_ADMIN"
  env:
    - "TZ=America/Chicago"
    - "VPN_ENABLED=yes"
    - "VPN_PROV=custom"
    - "VPN_CLIENT=openvpn"
    - "LAN_NETWORK=192.168.1.0/16"
```

### Field Meanings and Usage

*   **`your_service_identifier`**: The top-level key in your YAML file (e.g., `flaresolverr`, `deluge-vpn`). It acts as a unique identifier for this specific service configuration within the `local.services` map.
*   **`service_name`**: (Required) The human-readable name of the service, often used as the Docker container name.
*   **`image_name`**: (Required) The Docker image to pull, including its tag (e.g., `binhex/arch-delugevpn:2.2`).
*   **`host`**: (Optional) The server to deploy to. Must match a name in `config/servers.yaml`. Defaults to `tower`.
*   **`dns`**: DNS and proxy configuration for the service.
    *   **`dns.enabled`**: If `true`, Pulumi creates a DNS record and Nginx Proxy Manager host.
    *   **`dns.domain_name`**: The full domain name for the service (required if `dns.enabled` is true).
*   **`network`**: Network configuration for the service.
    *   **`network.internal`**: If `true`, the service is only accessible internally (not exposed via the internet).
    *   **`network.service_port`**: The port the service listens on *inside* the container.
    *   **`network.networks`**: A list of network objects to attach the container to. Each object has:
        *   `name`: The network name (e.g., `"br1"`).
        *   `ip_address`: (Optional) A static IP to assign on that network.
*   **`auth`**: Authentication configuration. See [stack-config.md](stack-config.md) for full auth options.
*   **`secrets`**: A map of environment variable names to Infisical secret names. Pulumi fetches these secrets from Infisical and injects them as environment variables.
*   **`env`**: (Optional) A list of `KEY=VALUE` strings that will be set as environment variables inside the container.
*   **`mounts`**: (Optional) A list of bind mount strings in the format `host_path:container_path[:ro]`.
*   **`privileged`**: (Optional) If `true`, runs the container in privileged mode (`--privileged`). Use only when capabilities or device passthrough are insufficient.
*   **`capabilities.add`**: (Optional) A list of Linux capabilities to add to the container (e.g., `CAP_NET_ADMIN`).
*   **`devices`**: (Optional) Device passthrough configuration.
    *   **`devices.gpu`**: If `true`, passes `/dev/dri` into the container (hardware transcoding).
    *   **`devices.usb`**: If `true`, passes `/dev/bus/usb` into the container (e.g. Coral TPU).
    *   **`devices.paths`**: A list of explicit `host_path:container_path` strings for other devices.
*   **`icon`**: (Optional) A URL to an icon for the service, often used by UI tools like Unraid.
