# Cloud Gateway Max

:::info
Last Updated 4/15/26
:::

## **1\. Device Overview**

- **Hostname/Device Name:** Cloud Gateway Max
- **Model:** Ubiquiti Cloud Gateway Max (UCG-Max-NS)
- **Role:** Primary Internet Router, Firewall & UniFi Network Controller
- **Management Access:**
    - **Primary IP Address for WebUI:** `192.168.1.1`
    - **WebUI Protocol:** HTTPS
    - **UniFi Controller:** Built-in (manages U7-Pro APs and USW-Pro-Max-16-PoE)

## **2\. Physical Information**

- **Physical Location:** Basement Rack
- **Serial Number:** [To be added later]
- **Power Backup (UPS):**
    - Connected to a UPS.
    - UPS Model: [To be added later]

## **3\. Network Configuration**

### **3.1. Port Layout**

| Port | Type | Connection |
|------|------|------------|
| 1 | RJ45 | Not Connected |
| 2 | RJ45 | Tower (MAC: `...7f:ee`) |
| 3 | RJ45 | [USW-Pro-Max-16-PoE](../Switches/USW-Pro-Max-16-PoE.md) |
| 4 | RJ45 | Not Connected |
| 5 | RJ45 | AT&T (WAN) |

### **3.2. WAN Interface**

- **Port 5 (WAN):**
    - **IP on this segment:** Assigned via DHCP from modem/gateway (`192.168.3.x/24`)
    - **Gateway for Internet Access:** `192.168.3.254`
    - **MTU:** `1500`
    - **Status:** Connected
    - **ISP:** AT&T
    - **Rated Speeds (Down/Up):** 1 Gbps / 1 Gbps

### **3.3. LAN Networks**

- **General Purpose — `192.168.1.0/24`:**
    - **Gateway IP:** `192.168.1.1`
    - **DHCP Server:** Enabled
        - **Range:** `192.168.1.2` – `192.168.1.254`
        - **DNS Servers provided to clients:** `192.168.4.53` (Technitium DNS)
    - **Status:** Active

- **Dev/Server — `192.168.4.0/23`:**
    - **Gateway IP:** `192.168.4.1`
    - **DHCP Server:** Enabled
        - **Range:** `192.168.4.2` – `192.168.5.254`
        - **DNS Servers provided to clients:** `192.168.4.53` (Technitium DNS)
    - **Status:** Active

- **WiFi Cameras — `192.168.8.0/28`:**
    - **Gateway IP:** `192.168.8.1`
    - **DHCP Server:** Enabled
        - **Range:** `192.168.8.2` – `192.168.8.14`
        - **DNS Servers provided to clients:** `192.168.4.53` (Technitium DNS)
    - **Status:** Work in Progress — no devices currently assigned

- **Planned — `192.168.2.0/26`:**
    - **Status:** Not in use — still being planned

### **3.4. Firewall & Port Forwarding**

- **General Policy Overview:**
    - Allow established/related connections.
    - Drop invalid state packets.
    - Default action for new incoming connections from WAN: Drop.

- **Port Forwarding Rules:**
    - **Plex:**
        - **External Port:** `32400`
        - **Protocol:** Both (TCP/UDP)
        - **Forward-to IP:** `192.168.1.41`
        - **Forward-to Port:** `32400`
    - **HTTP Web Server:**
        - **External Port:** `80`
        - **Protocol:** Both (TCP/UDP)
        - **Forward-to IP:** `192.168.4.2`
        - **Forward-to Port:** `80`
    - **HTTPS Web Server:**
        - **External Port:** `443`
        - **Protocol:** Both (TCP/UDP)
        - **Forward-to IP:** `192.168.4.2`
        - **Forward-to Port:** `443`

### **3.5. VPN Configuration**

- **Site-to-Site VPNs:** Not configured.
- **Remote Access VPN Server:** Not configured.

### **3.6. Quality of Service (QoS)**

- **Configured:** Not configured.

## **4\. Services & Features**

- **DHCP Server:** Enabled on all active LAN networks (see Section 3.3).
- **UniFi Network Controller:** Built-in. Manages:
    - [USW-Pro-Max-16-PoE](../Switches/USW-Pro-Max-16-PoE.md)
    - [U7-Pro (Basement)](../Access%20Points/U7-Pro%20Basement.md)
    - [U7-Pro (Living Room)](../Access%20Points/U7-Pro%20Living%20Room.md)

## **5\. Login & Access**

- **Admin Password Storage:** Stored in VaultWarden under "ui.com"

## **6\. Maintenance & Backup**

- **Configuration Backup:** No formal policy defined at this time.
- **Firmware Update Policy:** No formal policy defined at this time.

## **7\. Important Notes / Quirks**
