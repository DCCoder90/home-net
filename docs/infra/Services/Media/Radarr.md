# Radarr

:::info
Last Updated 3/18/2026
:::

## I. Service Overview

- **Service Name:** Radarr
- **Description:** Movie collection manager. Automatically searches for movies, manages downloads, and organizes media files.
- **Status:** Production
- **Managed By:** Pulumi IaC (`config/stacks/arr.yaml`)

## II. Docker & Network Configuration

- **Docker Image:** `linuxserver/radarr:6.0.4`
- **Network:** `br1`
- **IP Address:** `192.168.5.24`
- **Service Port:** `7878`
- **Domain:** `radarr.dcapi.app` (internal only)
- **Authentication:** Authentik proxy auth (group: `Arr`)

## III. IaC Configuration

- **Config File:** `config/stacks/arr.yaml` (part of the `arr_services` stack)
- **Environment Variables:**
    - `PUID=1000`
    - `PGID=100`
    - `TZ=America/Chicago`
- **Volume Mounts:**
    - `/mnt/user/Arr/radarr-data:/config`
    - `/mnt/user/Media:/media`
    - `/mnt/user/Downloads:/downloads`
    - `/etc/localtime:/etc/localtime:ro` (inherited from stack)
- **Secrets:** Proxy credentials (`radarr_username`, `radarr_password`) fetched from Infisical.

## IV. Application Configuration

- **Application Configuration:** Accessed via Web UI at `https://radarr.dcapi.app`.
    - Download client: Deluge
    - Indexers: Managed via Prowlarr

## V. Dependencies & Integrations

- **Internal Dependencies:**
    - [Prowlarr](Prowlarr.md) — Indexer management
    - [Deluge](../Other%20Services/Deluge.md) — Download client
- **Integrations:** Nginx Proxy Manager for reverse proxy, Authentik for SSO.

## VI. Useful Links & Notes

- **Official Docs:** https://wiki.servarr.com/radarr
- Ensure Deluge's completed download path is visible to Radarr under its `/downloads` mapping.
