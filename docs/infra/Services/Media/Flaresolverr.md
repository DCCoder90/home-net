# Flaresolverr

:::info
Last updated 5/27/25
:::

## **1\. Service Overview**

- **Service Name:** FlareSolverr
- **Docker Image Maintainer:** binhex
- **Purpose:** FlareSolverr runs a proxy server that utilizes Selenium with undetected-chromedriver to solve Cloudflare challenges. It returns the subsequent HTML and cookies, allowing other HTTP clients to bypass Cloudflare using these cookies. This is primarily used by other applications (like Prowlarr, Jackett, etc.) to access indexer sites protected by Cloudflare.
- **Status:** Production

---

## **2\. Host & Docker Details**

- **Running On:** [Tower](../../Hardware/Servers/Tower.md)
- **Unraid App Name / Container Name:** `binhex-flaresolverr`
- **Docker Image Repository:** `binhex/arch-flaresolverr`
- **Image Source Registry/Link:** https://github.com/binhex/arch-flaresolverr
- **Network Type:** Bridge

---

## **3\. Configuration & Data**

- **Volume Mounts (Data Persistence):**
    - **Container Path:** `/config`
    - **Host Path (Unraid):** `/mnt/user/appdata/binhex-flaresolverr` (Contains logs and potentially other configuration/cache)
    - **Container Path:** `/shared`
    - **Host Path (Unraid):** `binhex-shared`
- **Port Mapping:**
    - **Host Port:** `8191`
    - **Container Port:** `8191`
    - **Access URL (Local API):** `http://<TOWER_IP_OR_HOSTNAME>:8191` or `http://<FLARESOLVERR_CONTAINER_IP>:8191` (Primarily for other applications to use, not direct user access)
- **Key Environment Variables:**
    - `PUID=99` (User ID for permissions, typically 'nobody' on Unraid)
    - `PGID=100` (Group ID for permissions, typically 'users' on Unraid)
    - `UMASK=000` (File creation mask)
    - `SHARED_NETWORK=no` (Indicates it's not using a shared network stack like a VPN container, which is standard for bridge mode)
- **Backup Strategy:**
    - The `/mnt/user/appdata/binhex-flaresolverr` directory contains logs and configuration.

---

## **4\. Networking & Access**

- **Access URL (for client applications):** `http://<TOWER_IP_OR_HOSTNAME>:8191`
    - This URL is configured in applications like [Prowlarr](Prowlarr.md) to use FlareSolverr as a Cloudflare challenge solver.
- **IP Address:** Receives an IP address on a Docker bridge network. Outbound traffic NATted via "tower's" IP.

---

## **5\. Management & Troubleshooting**

- **Update Process:** Via Unraid Apps tab -> Check for Updates
- **Restart Procedure:** Via Unraid Docker UI.
- **Log Access:**
    - Via Unraid Docker UI (Logs button for the container).
    - Within `/mnt/user/appdata/binhex-flaresolverr/logs/` on the Unraid host.
- **Common Issues:**
    - Cloudflare challenges not solving (may require updating FlareSolverr, or Cloudflare has updated its protection).
    - High resource usage if many requests are processed simultaneously (though it's designed to be light when idle).

---

## **6\. Dependencies**

- [Tower](../../Hardware/Servers/Tower.md) **:** For Docker hosting.
- **Internet Connectivity:** To access target sites and solve challenges.

---

## **7\. Useful Links & Notes**

- **GitHub (Image Source):** https://github.com/binhex/arch-flaresolverr
- **FlareSolverr GitHub (Upstream Project):** https://github.com/FlareSolverr/FlareSolverr
- Primarily an intermediary service; not accessed directly by users via a web UI for general Browse.