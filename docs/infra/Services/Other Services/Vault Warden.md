# Vault Warden

:::info
Last updated 5/27/25
:::

## **1\. Service Overview**

- **Service Name:** Vaultwarden (Self-hosted Bitwarden server)
- **Purpose:** Provides a secure method for storing and retrieving passwords, passkeys, and other sensitive information.
- **Status:** Production

---

## **2\. Host & Docker Details**

- **Running On:** [Tower](../../Hardware/Servers/Tower.md)
- **Unraid App Name:** `vaultwarden`
- **Docker Image:** `vaultwarden/server`
- **Docker Hub Registry:** https://hub.docker.com/r/vaultwarden/server/
- **Network Type:** `br1` (Custom Docker network on "tower")
- **IP Address:** `192.168.4.201` (Static IP on `br1` network)
    - _This IP falls within the "Development" range of IP plan, needs to be updated._

---

## **3\. Configuration & Data**

- **Volume Mounts (Data Persistence):**
    - **Container Path:** `/data`
    - **Host Path (Unraid):** `/mnt/user/appdata/vaultwarden`
        - _This directory contains all Vaultwarden data, including user accounts, encrypted vaults, configuration, and attachments._ **_This is critical data._**
- **Port Mapping (Direct Container Access - primarily for initial setup or troubleshooting):**
    - **Host Port:** `4743`
    - **Container Port:** `80`
- **Key Environment Variables:**
    - `SIGNUPS_ALLOWED=false` (New user registrations via the web vault are disabled)
    - `INVITATIONS_ALLOWED=true` (Admins can invite new users; invitations are likely link-based as SMTP is not configured)
    - `WEBSOCKET_ENABLED=true` (Enables websocket notifications for client sync)
    - `ADMIN_TOKEN`: Set to the value stored for "Home Lab" in Vaultwarden (this token is used to access the `/admin` page). _Ensure this token is strong and securely managed._
- **SMTP Configuration:** SMTP is **not configured**.
    - Email-dependent features such as password reset via email or email notifications for new device logins will not function. Invitations will be link-based.
- **Backup Strategy:**
    - **CRITICAL NOTE:** There is **no defined backup strategy currently in place** for the `/mnt/user/appdata/vaultwarden` directory. This directory contains all user vaults, attachments, and configurations.
    - **Recommendation:** Implement a robust, regular, and tested backup solution for this directory immediately. Consider options like Unraid's Appdata Backup plugin, `rsync` to another machine/NAS, or a dedicated backup tool (Restic, Duplicati) to multiple locations (local and offsite/cloud). Losing this data could mean losing access to all stored credentials.

---

## **4\. Networking & Access**

- **Primary Access URL (Reverse Proxied):** `https://vault.dcapi.app` (Accessed via [Nginx-Proxy](../Networking/Nginx-Proxy.md) )
- **Local HTTP Access (primarily for troubleshooting/initial setup):** `http://192.168.4.201:4743`
- **Admin Page:** `https://vault.dcapi.app/admin` (Requires the `ADMIN_TOKEN` for access)

---

## **5\. Management & Troubleshooting**

- **Update Process:** Updated via Unraid's Apps tab / Update notifications (Action Center).
- **Restart Procedure:** Via Unraid Docker UI.
- **Log Access:** Via Unraid Docker UI.
- **Two-Factor Authentication (2FA):** Strongly recommended for all user accounts, especially admin accounts. No specific methods are enforced at an admin level; users can enable 2FA options individually within their accounts.
- **Emergency Access/Disaster Recovery:**
    - Securely store a copy of the `ADMIN_TOKEN`.
    - **Implement and regularly test backups of the** `/data` **volume.**
    - Document a restore procedure once backups are in place.

---

## **6\. Dependencies**

- [Tower](../../Hardware/Servers/Tower.md) :For Docker hosting.
- `br1` **Network:** Relies on this custom Docker network for its static IP.
- [Nginx-Proxy](../Networking/Nginx-Proxy.md) **:** For HTTPS termination and custom domain access.
- **Internet Connectivity:** For initial setup, pulling updates, and potentially for some 2FA methods or client access.

---

## **7\. Useful Links & Notes**

- **Official Vaultwarden Wiki:** [https://github.com/dani-garcia/vaultwarden/wiki](https://www.google.com/search?q=https://github.com/dani-garcia/vaultwarden/wiki)
- **Docker Hub:** https://hub.docker.com/r/vaultwarden/server/
- **Security Best Practices:**
    - Use strong, unique master passwords for all Vaultwarden user accounts.
    - All users should enable Two-Factor Authentication (2FA).
    - Keep the Vaultwarden container image and the underlying Unraid system updated.
    - **Regularly back up the** `/data` **directory (see critical note in section 3).**
    - Restrict access to the admin panel and protect the `ADMIN_TOKEN`.
    - Attachments are stored within the `/data` directory and are included in its backup (once implemented).