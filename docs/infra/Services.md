# Services

:::info
Last updated 3/18/2026
:::

# **Documentation Standard:**

Unless otherwise indicated on the specific service page, all services are hosted as Docker containers managed via Pulumi on the server known as [Tower](Hardware/Servers/Tower.md). Service definitions are maintained as YAML configuration files in the `config/stacks/` and `config/services/` directories, and are deployed automatically through the GitOps pipeline.

Documentation for services not following this standard will explicitly state their unique setup (e.g., cloud VPS, Raspberry Pi, dedicated VM).

# Core Services

These services are managed by the `core` Pulumi stack and must be applied before application services. They are not defined in `config/stacks/` or `config/services/` — instead they are provisioned directly in the `pulumi/` core program.

| Service | IP Address | Notes |
| --- | --- | --- |
| [Nginx-Proxy](Services/Networking/Nginx-Proxy.md) | 192.168.4.2 | Reverse proxy and SSL termination |
| [DNS-server](Services/Networking/DNS-server.md) | 192.168.4.53 | Technitium DNS server |

# Non-IaC Services

These services run on Tower but are not managed by Pulumi. They are configured manually via the Unraid Docker UI.

| Service | IP Address | Notes |
| --- | --- | --- |
| [Authentik](Services/Networking/Authentik.md) | 192.168.4.55 | SSO and identity provider |
| [Authentik-Worker](Services/Networking/Authentik-Worker.md) | 192.168.4.56 | Background worker for Authentik |
| [Vault Warden](Services/Other%20Services/Vault%20Warden.md) | 192.168.4.201 | Self-hosted Bitwarden password manager |
| [NetAlertX](Services/Networking/NetAlertX.md) | 192.168.1.45 | Network device monitoring |
| [OpenVPN-Client](Services/Networking/OpenVPN-Client.md) | Bridge | ProtonVPN exit node for Tailscale |

# IaC-Managed Services

Below is a list of all services managed by Pulumi via the `config/stacks/` and `config/services/` YAML definitions. These try to follow the IP allocation of the [IP Pool (tower reserved)](Home/IP%20Pool%20%28tower%20reserved%29.md).

## Stacks

### Arr Services (`config/stacks/arr.yaml`)

| Service | IP Address | Domain | Internal |
| --- | --- | --- | --- |
| Prowlarr | 192.168.5.22 | prowlarr.dcapi.app | Yes |
| Sonarr | 192.168.5.23 | sonarr.dcapi.app | Yes |
| [Radarr](Services/Media/Radarr.md) | 192.168.5.24 | radarr.dcapi.app | Yes |
| Lidarr | 192.168.5.25 | lidarr.dcapi.app | Yes |

### Media Servers (`config/stacks/media-servers.yaml`)

| Service | IP Address | Domain | Internal |
| --- | --- | --- | --- |
| Jellyfin | 192.168.5.46 | media.mallett.family | No |
| Tautulli | 192.168.5.30 | mediastats.dcapi.app | Yes |
| Plex | Host | — | No |

### Immich (`config/stacks/immich.yaml`)

| Service | IP Address | Domain | Internal |
| --- | --- | --- | --- |
| Immich Server | 192.168.5.47 | photos.mallett.family | No |
| Immich Postgres | Internal network | — | Yes |

### Karakeep (`config/stacks/karakeep.yaml`)

| Service | IP Address | Domain | Internal |
| --- | --- | --- | --- |
| Karakeep | 192.168.5.241 | keep.mallett.family | No |
| Meilisearch | Internal network | — | Yes |
| Chrome | Internal network | — | Yes |

### Monitoring (`config/stacks/monitoring.yaml`)

| Service | IP Address | Domain | Internal |
| --- | --- | --- | --- |
| Prometheus | 192.168.5.55 | — | Yes |
| Grafana | 192.168.5.56 | monitor.dcapi.app | Yes |
| What's Up Docker | 192.168.5.57 | docker.dcapi.app | Yes |

### Paperless (`config/stacks/paperless.yaml`)

| Service | IP Address | Domain | Internal |
| --- | --- | --- | --- |
| Paperless-NGX | 192.168.5.60 | paperless.mallett.family | No |
| Redis | Internal network | — | Yes |
| Postgres | Internal network | — | Yes |

### Temporal (`config/stacks/temporal.yaml`)

| Service | IP Address | Domain | Internal |
| --- | --- | --- | --- |
| Temporal Server | 192.168.5.62 | temporal-server.dcapi.app | Yes |
| Temporal UI | 192.168.5.61 | temporal.dcapi.app | Yes |
| Postgres | Internal network | — | Yes |

### AI Services (`config/stacks/ai.yaml`)

| Service | IP Address | Domain | Internal |
| --- | --- | --- | --- |
| n8n | 192.168.5.49 | n8n.dcapi.app | No |
| OpenClaw | 192.168.5.50 | claw.dcapi.app | Yes |
| [OpenClaw Voice](Services/AI/OpenClaw-Voice.md) | 192.168.5.50:3334 | voice.dcapi.app | No |

## Standalone Services

### Standalone Services (`config/services/`)

| Service | IP Address | Domain | Internal |
| --- | --- | --- | --- |
| [Requestrr](Services/Media/Requestrr.md) | 192.168.5.27 | — | Yes |
| [Deluge](Services/Other%20Services/Deluge.md) | 192.168.5.28 | — | Yes |
| [Flaresolverr](Services/Media/Flaresolverr.md) | 192.168.5.29 | — | Yes |
| Kavita | 192.168.5.33 | books.mallett.family | Yes |
| MySQL | 192.168.5.83 | — | Yes |
| Monica | 192.168.4.202 | — | Yes |
