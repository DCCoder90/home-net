# Radarr

:::info
Last Updated 5/27/25
:::

:::warning
Documentation incomplete or contains filler info. Needs to be completed, do not rely on.
:::

## I. Service Overview

- **Unraid App Name:** `binhex-radarr`
- **Description:** Movie collection manager.
- **Access URL(s):**
    - Unraid WebUI: Docker Tab -> radarr -> WebUI
    - Direct: `http://<UNRAID_IP>:7878` (if 7878 is host port)
    - Reverse Proxy: `https://radarr.yourdomain.com`
- **Status:** Production

## II. Unraid App & Docker Configuration

- **Source:** Community Applications (Binhex Repository)
- **Key Unraid Template Settings:**
    - **Host Path for AppData:** `/mnt/user/appdata/binhex-radarr`
    - **Host Path for /movies:** `/mnt/user/media/movies` (Container: `/movies`)
    - **Host Path for /downloads:** `/mnt/user/downloads/completed/movies` (Container: `/downloads`)
    - **WebUI Port:** `7878` (Container: `7878`)
    - **Network Type:** Bridge
- **Environment Variables (Notable):**
    - `UMASK=002`
    - `PUID=99` (nobody)
    - `PGID=100` (users)

## III. Application Configuration & Data

- **Application Configuration:** Accessed via Web UI.
    - Download client: Deluge (Unraid app, `http://<UNRAID_IP>:8112`)
    - Indexers: Managed via Prowlarr (Unraid app).
- **Credentials:**
    - Admin Username: `myradarradmin`
    - Password: Stored in Vaultwarden ("Radarr Admin").
- **Backup Strategy:** `appdata/binhex-radarr` backed up nightly by Unraid "Appdata Backup" plugin to local backup array and cloud.

## IV. Dependencies & Integrations

- **Internal Dependencies (Unraid Apps):**
    - [Prowlarr](link-to-prowlarr-page)
    - [Deluge](link-to-deluge-page)
- **Integrations:** [Nginx-Proxy-Manager App](link-to-npm-page) for SSL reverse proxy.

## V. Management & Troubleshooting (Unraid Context)

- **Update Process:** Via Unraid "Apps" tab -> "Check for Updates".
- **Restart Procedure:** Unraid UI -> Docker -> binhex-radarr -> Restart.
- **Log Access:** Unraid UI -> Docker -> binhex-radarr -> Logs.
    - App Data Logs: `/mnt/user/appdata/binhex-radarr/logs/`
- **Reinstall:** Use "Previous Apps" feature in Community Applications if needed, ensuring appdata path is preserved.

## VI. Useful Links & Notes

- **Support Thread:** \[Link to binhex-radarr Unraid support thread\]
- **Official Docs:** https://wiki.servarr.com/radarr
- **Notes:** Ensure Deluge's completed download path is visible to Radarr under its `/downloads` mapping.