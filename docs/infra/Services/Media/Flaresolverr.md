# Flaresolverr

:::info
Last updated 3/18/2026
:::

## **1\. Service Overview**

- **Service Name:** FlareSolverr
- **Purpose:** FlareSolverr runs a proxy server that utilizes Selenium with undetected-chromedriver to solve Cloudflare challenges. It returns the subsequent HTML and cookies, allowing other HTTP clients to bypass Cloudflare using these cookies. This is primarily used by other applications (like Prowlarr) to access indexer sites protected by Cloudflare.
- **Status:** Production
- **Managed By:** Terraform IaC (`config/services/flaresolverr.yaml`)

---

## **2\. Host & Docker Details**

- **Running On:** [Tower](../../Hardware/Servers/Tower.md)
- **Docker Image:** `flaresolverr/flaresolverr:v3.4.6`
- **Network Type:** `br1` (Static IP)
- **IP Address:** `192.168.5.29`
- **Service Port:** `8112`

---

## **3\. Configuration**

- **IaC Config File:** `config/services/flaresolverr.yaml`
- **DNS:** Disabled (internal service, accessed by other containers directly)
- **Authentication:** Disabled
- **Access URL (for client applications):** `http://192.168.5.29:8112`
    - This URL is configured in applications like [Prowlarr](Prowlarr.md) to use FlareSolverr as a Cloudflare challenge solver.

---

## **4\. Dependencies**

- [Tower](../../Hardware/Servers/Tower.md): For Docker hosting.
- **Internet Connectivity:** To access target sites and solve challenges.

---

## **5\. Useful Links & Notes**

- **FlareSolverr GitHub (Upstream Project):** https://github.com/FlareSolverr/FlareSolverr
- Primarily an intermediary service; not accessed directly by users via a web UI.
