# Services

:::info
Last updated 7/8/2025
:::

# **Documentation Standard:**

Unless otherwise indicated on the specific service page, all services are hosted as Docker containers managed by Unraid on the server known as [Tower](Hardware/Servers/Tower.md). Configuration typically involves settings within the Unraid Docker template and the application's own web interface. Documentation for services not following this standard will explicitly state their unique setup (e.g., cloud VPS, Raspberry Pi, dedicated VM).

# Services

Below is a list of services along with their IP address that are available within the HomeLab network. These try to follow the IP allocation of the [IP Pool (tower reserved)](Home/IP%20Pool%20%28tower%20reserved%29.md).

All addresses listed as “Bridge” will be allocated to the IP address for eth1 on [Tower](Hardware/Servers/Tower.md). This is 192.168.4.41 at the time of this writing.

| Service | IP Address |
| --- | --- |
| [Authentik](Services/Networking/Authentik.md) | 192.168.4.55 |
| [Authentik-Worker](Services/Networking/Authentik-Worker.md) | 192.168.4.56 |
| [Radarr](Services/Media/Radarr.md) | 192.168.5.24 |
| [Sonarr](Services/Media/Sonarr.md) | 192.168.5.23 |
| [Prowlarr](Services/Media/Prowlarr.md) | 192.168.5.22 |
| [Lidarr](Services/Media/Lidarr.md) | 192.168.5.25 |
| [Plex](Services/Media/Plex.md) | Host |
| [Requestrr](Services/Media/Requestrr.md) | 192.168.5.27 |
| [Postgres14](Services/Database/Postgres14.md) | 192.168.4.105 |
| [Redis](Services/Database/Redis.md) | 192.168.4.115 |
| [Nginx-Proxy](Services/Networking/Nginx-Proxy.md) | 192.168.4.2 |
| [Flaresolverr](Services/Media/Flaresolverr.md) | 192.168.5.29 |
| [Deluge](Services/Other%20Services/Deluge.md) | 192.168.5.28 |
| [Jellyfin](Services/Media/Jellyfin.md) | 192.168.4.46 |
| [DNS-server](Services/Networking/DNS-server.md) | 192.168.1.59 |
| [Immich](Services/Other%20Services/Immich.md) | 192.168.5.47 |
| [MongoDB](Services/Database/MongoDB.md) | 192.168.5.74 |
| [NetAlertX](Services/Networking/NetAlertX.md) | 192.168.1.45 |
| [OpenVPN-Client](Services/Networking/OpenVPN-Client.md) | Bridge |
| [PostgresSQL_Immich](Services/Database/PostgresSQL_Immich.md) | 192.168.5.48 |
| [Vault Warden](Services/Other%20Services/Vault%20Warden.md) | 192.168.4.201 |