# Nginx-Proxy

:::info
Last Updated 5/27/25
:::

:::warning
See:  
4\. Networking & Access need to update with correct information  
5\. Management & Troubleshooting need to update with correct information
:::

:::danger
No backup strategy determined. Need to establish clear backup standards.
:::

## **1\. Service Overview**

- **Service Name:** Nginx Proxy Manager
- **Purpose:** Provides an easy-to-use interface for managing Nginx proxy hosts, including SSL certificate generation and management via Let's Encrypt. It simplifies exposing internal services to the internet with custom domains and HTTPS.
- **Status:** Production

---

## **2\. Host & Docker Details**

- **Running On:** [Tower](../../Hardware/Servers/Tower.md)
- **Unraid App Name / Container Name:** `Nginx-Proxy-Manager-Official`
- **Docker Image Repository:** `jc21/nginx-proxy-manager`
- **Docker Hub Registry:** https://hub.docker.com/r/jc21/nginx-proxy-manager
- **Network Type:** `br1` (Custom Docker network on "tower")
- **IP Address:** `192.168.4.2` (Static IP on `br1` network)
- **Extra Parameters:** `--memory=1G --no-healthcheck`

---

## **3\. Configuration & Data**

- **Volume Mounts (Data Persistence):**
    - **Container Path:** `/data`
        - **Host Path (Unraid):** `/mnt/user/appdata/NginxProxymanager-Latest/data/` (Contains SQLite database, proxy configurations, access lists, etc. **Critical to back up.**)
    - **Container Path:** `/etc/letsencrypt`
        - **Host Path (Unraid):** `/mnt/user/appdata/NginxProxymanager-Latest/letsencrypt/` (Contains Let's Encrypt SSL certificates and account information. **Critical to back up.**)
    - **Container Path:** `/var/log`
        - **Host Path (Unraid):** `/tmp/NginxProxymanager-Latest/var/log`
- **Port Mappings:**
    - Ports that require external access are mapped accordingly on [EdgeRouter Pro](../../Hardware/Routers/EdgeRouter%20Pro.md) . Typically, these will only be 80 and 443.
    - **Admin WebUI:**
        - **Host Port:** `81`
        - **Container Port:** `81`
    - **Proxy HTTP Traffic:**
        - **Host Port:** `80`
        - **Container Port:** `80`
    - **Proxy HTTPS Traffic:**
        - **Host Port:** `443`
        - **Container Port:** `443`
    - **Application Websocket/API (Potentially for internal NPM use, not typically exposed):**
        - **Host Port:** `3000`
        - **Container Port:** `3000`
- **Key Environment Variables:**
    - `DB_SQLITE_FILE=/data/database.sqlite`
    - `PUID` / `PGID`: Not explicitly set by user; container runs with image default user permissions.
- **Backup Strategy:**
    - **CRITICAL NOTE:** There is **no current backup strategy defined** for the critical data directories:
        - `/mnt/user/appdata/NginxProxymanager-Latest/data/`
        - `/mnt/user/appdata/NginxProxymanager-Latest/letsencrypt/`
    - **Recommendation:** Implement a robust, regular, and tested backup solution for these directories immediately. Losing this data would require reconfiguring all proxy hosts and re-issuing SSL certificates. Consider Unraid's Appdata Backup plugin, `rsync` to another machine/NAS, or a dedicated backup tool.

---

## **4\. Networking & Access**

- **Admin WebUI:** `http://192.168.4.2:81`
- **Firewall Considerations:**
    - [EdgeRouter Pro](../../Hardware/Routers/EdgeRouter%20Pro.md) is configured to forward external ports `80` (for HTTP traffic and HTTP-01 Let's Encrypt challenges) and `443` (for HTTPS traffic) to `192.168.4.2` (this Nginx Proxy Manager container

---

## **5\. Management & Troubleshooting**

- **Admin Credentials (NPM WebUI):**
    - **Username/Email:** Stored in Vaultwarden under "Home Lab"
    - **Password:** Stored in Vaultwarden under "Home Lab"
- **Update Process:** Updated via Unraid's Apps tab / Docker update notifications.
- **Restart Procedure:** Via Unraid Docker UI.
- **Log Access:**
    - Via Unraid Docker UI (for the container's stdout/stderr).
    - Nginx access and error logs are mapped to `/tmp/NginxProxymanager-Latest/var/log` on the Unraid host (note: `/tmp` is typically not persistent across host reboots).
- **SSL Certificate Renewal:** Handled automatically by Nginx Proxy Manager using Let's Encrypt (requires port 80 to be externally accessible for HTTP-01 challenges, or DNS challenge configuration if used).
    - There is a certificate for \*.dcapi.app that is registered through cloudflare that needs to be updated manually. This needs to be noted.

---

## **6\. Dependencies**

- [Tower](../../Hardware/Servers/Tower.md) **:** For Docker hosting.
- [EdgeRouter Pro](../../Hardware/Routers/EdgeRouter%20Pro.md) **:** For external port forwarding (80/443) to this service.
- **Let's Encrypt:** For SSL certificate generation and renewal.
- [Cloudflare](../External%20Services/Cloudflare.md) : For domain handling and SSL certificate generation and renewal

---

## **7\. Useful Links & Notes**

- **Official Nginx Proxy Manager Website:** https://nginxproxymanager.com/
- **Docker Hub:** https://hub.docker.com/r/jc21/nginx-proxy-manager
- The log path `/tmp/NginxProxymanager-Latest/var/log` on Unraid means Nginx-specific logs will likely be lost on host reboot. If persistent Nginx logs are required for extended debugging, consider remapping the container's `/var/log` to a persistent path within your appdata share (e.g., `/mnt/user/appdata/NginxProxymanager-Latest/log/nginx/`).