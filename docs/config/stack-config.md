
## ⚙️ Stack Configuration Files

Stack configuration files define groups of related services and their specific settings. These files are written in YAML and are automatically discovered and loaded by Pulumi.

### File Location

All stack configuration files must be placed in the `config/stacks/` directory. Each file should have a `.yaml` extension. Pulumi will merge all `.yaml` files found in this directory into a single configuration object.

### Structure

Each stack configuration file should define a top-level key that represents the name of the stack (e.g., `karakeep`, `arr_services`). Under this key, you define global settings for the stack and a list of `services`.

```yaml
your_stack_name:
  # Optional: Server to deploy all services in this stack to. Defaults to "tower".
  # Must match a name defined in config/servers.yaml.
  host: "myserver"
  # Optional: Global environment variables for all services in this stack
  env:
    - "GLOBAL_ENV_VAR=value"
  # Optional: Global bind mounts for all services in this stack
  mounts:
    - "/path/on/host:/path/in/container"
  # Optional: Global Docker volumes for all services in this stack
  volumes:
    - "my_volume:/data"
  # Optional: List of secret names to auto-generate as stable random values.
  # On first deploy, each name is created as a secret in Infisical under /generated.
  # The value is then injected into every service in this stack as KEY=value.
  generated_secrets:
    - "API_KEY"
    - "DATABASE_PASSWORD"
  # Optional: Custom Docker networks to be created for this stack
  networks:
    my_custom_network:
      internal: true
      driver: "bridge"
  # Required: A map of services to be deployed within this stack
  services:
    service_one_key: # Unique key for the service within this stack
      # Required: Display name of the service
      service_name: "My Awesome Service"
      # Optional: Override the stack-level host for this specific service.
      host: "myserver"
      # Required: Docker image name and tag
      image_name: "myrepo/my-service:latest"
      # Optional: Docker network mode (e.g., "host", "bridge")
      network_mode: "bridge"
      # Optional: Docker restart policy. Default: "unless-stopped"
      # Valid values: no, always, on-failure, unless-stopped
      restart_policy: "unless-stopped"
      # Optional: Description of the service
      description: "A brief description of what this service does."
      # Optional: URL to an icon for the service (e.g., for Unraid UI)
      icon: "https://example.com/icon.png"
      # Optional: Full domain name for the service (used by some modules alongside dns.domain_name)
      domain_name: "service.yourdomain.com"
      # Optional: Device passthrough for the container
      devices:
        gpu: true         # Passes /dev/dri (GPU/hardware transcoding)
        usb: true         # Passes /dev/bus/usb (e.g. Coral TPU)
        paths:            # Explicit host:container device paths
          - "/dev/ttyUSB0:/dev/ttyUSB0"
      # Optional: Service-specific environment variables
      env:
        - "SERVICE_ENV_VAR=value"
        - "SECRET_VAR=${API_KEY}" # References a generated secret
      # Optional: Service-specific bind mounts
      mounts:
        - "/path/on/host:/path/in/container"
      # Optional: Service-specific Docker volumes
      volumes: [] # List of volume objects (see docker module variables for structure)
      # Optional: List of commands to run in the container.
      commands: ["--some-flag", "value"]
      # Optional: Run the container in privileged mode
      privileged: true
      # Optional: Linux capabilities to add/drop for the container
      capabilities:
        add: ["NET_ADMIN"]
        drop: ["SYS_ADMIN"]
      # Optional: Secrets fetched from Infisical and injected as environment variables.
      # Keys are env var names, values are Infisical secret names.
      secrets:
        MY_API_KEY: "infisical_secret_name"
        DB_PASSWORD: "infisical_db_password"
      # Optional: Application config files to write to the remote host before the container starts.
      # Each entry fetches the Infisical key from the /config folder and writes it to the host path.
      # The file is only re-written when its content changes. Add a corresponding entry to
      # mounts: to make the file available inside the container.
      configfiles:
        - path: "/mnt/user/appdata/myapp/config.yaml"  # Absolute path on the remote host
          key: "MYAPP_CONFIG"                           # Infisical key in /config folder
          permissions: "0600"                           # Optional; defaults to "0600"
      # Optional: Network configuration for the service
      network:
        # Optional: Is this service internal-only (not accessible from the internet)? Default: false
        internal: false
        # Optional: The port the service listens on inside the container
        service_port: 8080
        # Optional: List of networks to attach the container to.
        # Each entry is an object with a name and optional ip_address.
        networks:
          - name: "br1"
            ip_address: "192.168.5.100"
          - name: "my_custom_network"
      # Optional: DNS configuration for the service
      dns:
        # Optional: Enable DNS record and proxy host creation. Default: false
        enabled: true
        # Optional: The full domain name for the service. Required if dns.enabled is true.
        domain_name: "service.yourdomain.com"
      # Optional: Authentication configuration for the service
      auth:
        # Optional: Enable authentication for this service. Default: false
        enabled: true
        # Optional: The primary Authentik group — members get full/admin access. Default: service name
        group: "My App Group"
        # Optional: Additional Authentik groups — each gets its own policy binding to the application.
        # Useful for role-based access (e.g. a viewer group with read-only access).
        additional_groups:
          - "My App Viewers"
        # Optional: Enable proxy authentication (e.g., Authentik forward auth).
        proxy:
          enabled: true
          user_secret: "service_username"  # Infisical secret name for proxy username
          pass_secret: "service_password"  # Infisical secret name for proxy password
          auth_secret_name: "infisical_key"  # Optional: Infisical key for X-Proxy-Secret header value
        # Optional: OAuth configuration for the service
        oauth:
          # Optional: Enable OAuth authentication. Default: false
          enabled: true
          # Optional: Map of environment variable names to Authentik OAuth output keys.
          keys:
            OAUTH_CLIENT_ID: "client_id"
            OAUTH_CLIENT_SECRET: "client_secret"
            OAUTH_WELL_KNOWN_URL: "provider_info_url"
          # Optional: List of OAuth scopes to request.
          scopes:
            - "openid"
            - "profile"
            - "email"
          # Optional: Additional redirect URIs for the OAuth provider.
          # These are appended to the base domain redirect URI.
          redirect_uris:
            - "/oauth/callback"
            - "/login/oauth/callback"

    service_two_key: # Another service definition...
      # ...
```

### Field Meanings and Usage

*   **`host` (Stack/Service Level)**: The server to deploy to. Must match a name in `config/servers.yaml`. Defaults to `tower`. Set at the stack level to deploy all services in the stack to one host, or per-service to override. See [Adding a new server](../setup.md#9-adding-a-new-server).
*   **`env` (Stack/Service Level)**: A list of `KEY=VALUE` strings for environment variables. Service-level `env` is merged with stack-level `env`. Values like `${SECRET_NAME}` will be replaced by the corresponding `generated_secrets` value at deploy time.
*   **`mounts` (Stack/Service Level)**: A list of bind mount strings in `host_path:container_path[:ro]`. Service-level `mounts` are merged with stack-level `mounts`.
*   **`commands` (Service Level)**: A list of strings representing the command to run in the container, overriding the image's default command.
*   **`volumes` (Stack/Service Level)**: A list of Docker volume configurations.
*   **`generated_secrets`**: A list of string names (e.g., `"DB_PASSWORD"`) for secrets that are auto-generated on first deploy and stored in Infisical under `/generated`. Each secret is automatically injected into every service in the stack as `KEY=value`. You can also reference a generated secret inside another env var using `${SECRET_NAME}` syntax (e.g., `"DB_PASS=${DB_PASSWORD}"`) if you need the env var name to differ from the key name. Values persist across runs, making them safe for stateful services like databases.
*   **`networks` (Stack Level)**: Defines custom Docker networks to be created for this stack. These are separate from `br0` and `br1`.
*   **`services`**: The core of the stack, defining individual Docker containers.
    *   **`service_name`**: The name of the Docker container and the base for Authentik application names.
    *   **`image_name`**: The Docker image to pull (e.g., `linuxserver/sonarr:4.0.16`).
    *   **`network_mode`**: The Docker network mode to set (e.g., `"host"` for services like Plex that need host networking).
    *   **`restart_policy`**: Docker restart policy for the container. Valid values: `no`, `always`, `on-failure`, `unless-stopped`. Defaults to `unless-stopped`.
    *   **`privileged`**: If `true`, runs the container in privileged mode (`--privileged`). Grants full host device access and disables most security isolation — use only when capabilities or device passthrough are insufficient.
    *   **`devices`**: Controls device passthrough into the container.
        *   **`devices.gpu`**: If `true`, passes `/dev/dri` into the container (used for hardware transcoding in media servers).
        *   **`devices.usb`**: If `true`, passes `/dev/bus/usb` into the container (used for USB accelerators such as a Coral TPU).
        *   **`devices.paths`**: A list of explicit `host_path:container_path` device strings for any other device passthrough needs.
    *   **`secrets`**: A map of environment variable names to Infisical secret names. Unlike `generated_secrets` (which are stack-level and referenced via `${}`), these are per-service and directly injected as environment variables.
    *   **`configfiles`**: A list of application config files to write to the remote host before the container starts. Each entry has:
        *   `path`: The absolute path on the remote host where the file will be written (e.g., `/mnt/user/appdata/myapp/config.yaml`).
        *   `key`: The Infisical secret key in the `/config` folder whose value becomes the file's content.
        *   `permissions`: (Optional) Octal permission string for the file (e.g., `"0644"`). Defaults to `"0600"`.
        The file is only re-written when its content changes between deploys. Running `pulumi destroy` removes the file. Add a corresponding `mounts:` entry to make the file available inside the container. Fails hard if the Infisical key is not found in `/config`.
    *   **`network.internal`**: If `true`, the service is not exposed externally via Nginx Proxy Manager.
    *   **`network.service_port`**: The port the service listens on *inside* the container.
    *   **`network.networks`**: A list of network objects to attach the container to. Each object has:
        *   `name`: The network name (e.g., `"br1"`, `"my_custom_network"`).
        *   `ip_address`: (Optional) A static IP address to assign on that network.

    > **Guidance**:
    > - Use `br0` for services that need to be directly accessible on your primary LAN (e.g., for network discovery protocols like mDNS or for devices that need to connect directly).
    > - Use `br1` for most standard services that will be accessed via the reverse proxy. This isolates their traffic from the main LAN.

    *   **`dns.enabled`**: If `true`, Pulumi will create a DNS record and an Nginx Proxy Manager host for this service.
    *   **`dns.domain_name`**: The full domain name for the service (e.g., `sonarr.dcapi.app`).
    *   **`dns.access_list_id`**: This field is **not set directly**. The module automatically assigns the "Internal Only" access list for internal services and the "CloudFlare Only" access list for external services.
    *   **`auth.enabled`**: If `true`, authentication is configured for this service in Authentik (group, application, policy binding).
    *   **`auth.group`**: The primary Authentik group for this application. Members of this group are granted access. Multiple services can share a group. Defaults to the service name.
    *   **`auth.additional_groups`**: A list of extra Authentik group names. Pulumi creates each group and adds a policy binding for it to the application. Useful for role-based access — e.g. a viewer group alongside an admin group. Each service consuming the `x-forwarded-groups` header (like Frigate's `role_map`) can then map these groups to application-level roles.
    *   **`auth.proxy`**: Enables Authentik forward-auth proxy for this service. When enabled, NPM routes all traffic for this service's domain through the Authentik proxy outpost instead of directly to the service. The outpost checks authentication and forwards authenticated requests on to the service. Requires the three-phase Authentik bootstrap — see [Authentik bootstrap](../setup.md#authentik-bootstrap).
    *   **`auth.proxy.auth_secret_name`**: Optional. The Infisical secret key whose value is injected as an `X-Proxy-Secret` header by NPM. The upstream service must validate this header to ensure requests arrive via the proxy. Generate the secret with `openssl rand -hex 32` and store it in Infisical.
    *   **`auth.oauth`**: Configures an Authentik OAuth2/OIDC provider for this service. The client ID, client secret, and issuer URL are exported as Pulumi stack outputs and can be read by the service via `auth.oauth.keys`.
    *   **`auth.oauth.enabled`**: If `true`, an Authentik OAuth2 provider and application are created for this service.
    *   **`auth.oauth.keys`**: Maps desired environment variable names (e.g., `OAUTH_CLIENT_ID`) to Authentik OAuth provider output attributes. Common values are `client_id`, `client_secret`, and `provider_info_url`.
    *   **`auth.oauth.scopes`**: A list of OAuth scopes to request from Authentik (e.g., `openid`, `profile`, `email`).
    *   **`auth.oauth.redirect_uris`**: A list of additional relative paths (e.g., `/oauth/callback`) that will be appended to the service's domain name to form the complete, valid OAuth redirect URIs required by Authentik.

---

## Proxy authentication outpost

When one or more services use `auth.proxy`, Pulumi manages a dedicated Authentik proxy outpost. The outpost is a separate container (`ghcr.io/goauthentik/proxy`) that sits between NPM and your services, handling authentication before forwarding traffic.

### How it works

```
Browser → NPM → Authentik Proxy Outpost → Service
                      ↑
               checks auth with Authentik
```

1. NPM receives a request for a protected domain
2. NPM forwards the request to the outpost (instead of directly to the service)
3. The outpost checks whether the user is authenticated with Authentik
4. If not authenticated, the user is redirected to the Authentik login page
5. If authenticated, the outpost proxies the request to the actual service

### Configuration in `system.yaml`

Add the following to the `authentik` block:

```yaml
authentik:
  ip_address: "192.168.4.55"
  port: 9000
  domain_name: "auth.example.com"
  # ... existing fields ...

  # Dedicated proxy outpost
  outpost_ip_address: "192.168.4.56"  # an unused IP on br1 — required
  outpost_port: 9000                   # optional, defaults to 9000
  outpost_image: "ghcr.io/goauthentik/proxy:2024.10.x"  # optional, defaults to :latest
```

- **`outpost_ip_address`** — A static IP on `br1` for the outpost container. Must be set to enable the outpost. NPM will forward proxy-auth traffic here.
- **`outpost_port`** — The port the outpost listens on. Defaults to `9000`.
- **`outpost_image`** — The outpost container image. Pin this to match your Authentik server version (e.g., `ghcr.io/goauthentik/proxy:2024.10.3`). Defaults to `:latest`.

### Setup

See [Authentik bootstrap](../setup.md#authentik-bootstrap) for the full three-phase deploy sequence. In summary:

1. Set `outpost_ip_address` in `system.yaml` and add a service with `auth.proxy.enabled: true`
2. Run `pulumi up` twice (phases 1 and 2) to deploy Authentik and create the outpost record
3. In the Authentik UI: **Applications → Outposts → "Pulumi Proxy Outpost" → View Deployment Info** — copy the token
4. Store the token in Infisical as `authentik_outpost_token`
5. Run `pulumi up` (phase 3) to deploy the outpost container

After the initial bootstrap, the outpost record is fully Pulumi-managed — proxy providers for new services are added to the outpost automatically on each `pulumi up`. No manual assignment in the Authentik UI is needed.
