# EdgeRouter Pro

:::info
Last Updated 5/27/25
:::

## **1\. Device Overview**

- **Hostname/Device Name:** EdgeRouter-Pro-8-Port
- **Model:** Ubiquiti EdgeRouter Pro 8-Port (ERP-8)
- **Role:** Primary Internet Router & Firewall
- **Firmware Version:** EdgeRouter Pro v2.0.9-hotfix.7
- **Management Access:**
    - **Primary IP Address for WebUI:** `192.168.1.1` (via `eth1`)
    - **WebUI Protocol:** HTTPS
    - **SSH Access:** Disabled

## **2\. Physical Information**

- **Physical Location:** Basement Rack
- **Serial Number:** _(User may add if desired)_
- **Power Backup (UPS):**
    - Connected to a UPS.
    - UPS Model: _(User may add if desired)_

## **3\. Network Configuration**

### **3.1. WAN Interface(s)**

- **Interface** `eth0` **(Label: Internet):**
    - **Router's IP on this segment:** `192.168.3.1/24` (Assigned via DHCP from modem/gateway)
    - **Gateway for Internet Access:** `192.168.3.254` (static route `0.0.0.0/0 via 192.168.3.254 dev eth0` is active in FIB)
    - **IPv6:** `dhcpv6` is configured on this interface, but it does not obtain an IPv6 address/prefix; IPv6 is not actively used.
    - **MTU:** `1500`
    - **Status:** Connected
    - **ISP (Internet Service Provider):** AT&T
    - **Modem/ONT Connected To:** Refer to separate "Modem Documentation Page" for details.
    - **Actual Public IP Address (as seen by internet):** `172.5.27.161` (Provided by AT&T upstream of the device at `192.168.3.254`)
    - **Rated Speeds (Down/Up):** 1 Gbps / 1 Gbps

### **3.2. LAN Interface(s) / VLANs / Bridged Interfaces**

- **Interface** `eth1` **(Label: Wired, Corresponds to LAN1):**
    - **Network Name/Purpose:** Wired LAN (Primary management LAN, Plex Server LAN)
    - **IP Subnet & Router IP:** `192.168.1.0/24`, Router IP: `192.168.1.1`
    - **MTU:** `1500`
    - **Status:** Connected
    - **DHCP Server:** Enabled.
        - **Range Start:** `192.168.1.38`
        - **Range Stop:** `192.168.1.191`
        - **Lease Time:** `86400` seconds (24 hours)
        - **DNS Servers provided to clients:** `192.168.1.1` (Router itself)
- **Interface** `eth2` **(Label: WiFi, Corresponds to LAN2):**
    - **Network Name/Purpose:** WiFi Network LAN
    - **IP Subnet & Router IP:** `192.168.2.0/24`, Router IP: `192.168.2.1`
    - **MTU:** `1500`
    - **Status:** Connected
    - **DHCP Server:** Enabled.
        - **Range Start:** `192.168.2.38`
        - **Range Stop:** `192.168.2.243`
        - **Lease Time:** `86400` seconds (24 hours)
        - **DNS Servers provided to clients:** `8.8.8.8`, `1.1.1.1`
- **Interface** `eth3` **(Label: eth3):**
    - **Network Name/Purpose:** Unused.
    - **IP Subnet & Router IP:** `192.168.9.0/24`, Router IP: `192.168.9.1`
    - **MTU:** `1500`
    - **Status:** Disconnected
    - **DHCP Server:** Not configured.
- **Interface** `eth4` **(Label: Dev, Corresponds to LAN4):**
    - **Network Name/Purpose:** Physical connection to eth1 on [Tower](../Servers/Tower.md)
    - **IP Subnet & Router IP:** `192.168.4.0/23`, Router IP: `192.168.4.1`
    - **MTU:** `1500`
    - **Status:** Connected
    - **DHCP Server:** Enabled.
        - **Range Start:** `192.168.4.38`
        - **Range Stop:** `192.168.5.254`
        - **Lease Time:** `86400` seconds (24 hours)
        - **DNS Servers provided to clients:** `1.1.1.1`, `8.8.8.8`
- **Interfaces** `eth5`**,** `eth6`**,** `eth7`**:**
    - **Status:** Disconnected, not configured, and serve no current purpose.
    - **MTU:** `1500`
- **Bridged Interfaces (e.g.,** `switch0`**):** No software bridges are configured.

### **3.3. Firewall**

- **Port Forwarding Configuration:**
    - **WAN Interface used for Port Forwarding:** `eth0`
    - **LAN Interface(s) targeted by Port Forwarding:** `eth1` (Wired LAN), `eth4` (Dev Network)
    - **Hairpin NAT:** Enabled
- **General Policy Overview:**
    - Allow established/related connections.
    - Drop invalid state packets.
    - Default action for new incoming connections from WAN: Drop.
- **Key Custom Firewall Rule Sets/Groups & Specific Inbound Allow Rules:**
    - **Twilio SIP Connectivity (Applied to** `WAN_IN` **ruleset):**
        - **Source IP Address(es):** `168.86.128.0/18` (Covers `168.86.128.0` - `168.86.191.255`)
        - **Destination Port(s):** `10000-60000`
        - **Protocol:** UDP
        - **Purpose:** Allow SIP connectivity from Twilio services.
        - **Destination:** Currently "empty" / "allowed to all." _(User Note: This configuration is acknowledged and needs to be updated/reviewed for appropriate scoping eventually. Current state implies traffic is allowed to the router's WAN IP on these ports, or potentially any internal IP if not further restricted by DNAT rules. Verification of the intended destination and scope is recommended for security.)_
    - No other named custom firewall rule sets are configured.
- **Actual Port Forwarding Rules:**
    - **Plex Custom Port 1:**
        - **Description:** Plex (Custom Port 33443)
        - **Original (External) Port:** `33443`
        - **Protocol:** Both (TCP/UDP)
        - **Forward-to (Internal) IP Address:** `192.168.1.41` (on `eth1`)
        - **Forward-to (Internal) Port:** `33443`
    - **Plex Main Port (via custom external):**
        - **Description:** Plex (Main via 55001)
        - **Original (External) Port:** `55001`
        - **Protocol:** Both (TCP/UDP)
        - **Forward-to (Internal) IP Address:** `192.168.1.41` (on `eth1`)
        - **Forward-to (Internal) Port:** `32400`
    - **HTTP Web Server:**
        - **Description:** HTTP Web Server
        - **Original (External) Port:** `80`
        - **Protocol:** Both (TCP/UDP)
        - **Forward-to (Internal) IP Address:** `192.168.4.2` (on `eth4`)
        - **Forward-to (Internal) Port:** `80`
    - **HTTPS Web Server:**
        - **Description:** HTTPS Web Server
        - **Original (External) Port:** `443`
        - **Protocol:** Both (TCP/UDP)
        - **Forward-to (Internal) IP Address:** `192.168.4.2` (on `eth4`)
        - **Forward-to (Internal) Port:** `443`
- **NAT Exclusions (if any):** None identified.

### **3.4. VPN Configuration**

- **Site-to-Site VPNs:** Not configured.
- **Remote Access VPN Server:** Not configured.

### **3.5. Quality of Service (QoS)**

- **Configured:** Not configured.

### **3.6. DNS**

- **LAN DNS Forwarding:**
    - Router acts as a DNS forwarder for clients on `eth1` (Wired LAN) using `192.168.1.1`.
    - Router acts as a DNS forwarder for clients on `eth2` (WiFi LAN) using `192.168.2.1`.
    - Clients on `eth4` (Dev LAN) use DNS servers `1.1.1.1` and `8.8.8.8` directly as provided by DHCP and do not use `192.168.4.1` for DNS forwarding.
- **Upstream DNS Servers used by Router:** Specific upstream DNS servers are not manually set. The router uses DNS servers provided by the ISP (AT&T) via DHCP on the `eth0` (WAN) interface for its forwarding services for `eth1` and `eth2`. _(User Note: The DNS forwarding setup or the reliance on ISP DNS is under review and may be changed/disabled.)_
- **Conditional Forwarding:** No conditional forwarding is configured.
- **DNS Host Records / Static DNS entries:** No static host entries are configured.

### **3.7. Static Routes (beyond default and connected)**

- **Configured (Inactive):**
    - `0.0.0.0/0` via `192.168.100.1` (Not in FIB) - Purpose unknown/legacy.
    - `0.0.0.0/0` via `173.18.93.1` (Not in FIB) - Purpose unknown/legacy.
- No other custom static routes are configured.

### **3.8. Dynamic Routing**

- No dynamic routing protocols are in use.

## **4\. Services & Features Used**

- **Hardware Offloading:** Not enabled (all specific offload functions are disabled).
- **UPnP / UPnP2:** Not enabled/Not set up.
- **DPI (Deep Packet Inspection):** Enabled.
- **NTP Client:** Enabled. Router syncs time from:
    - `0.ubnt.pool.ntp.org`
    - `1.ubnt.pool.ntp.org`
    - `2.ubnt.pool.ntp.org`
    - `3.ubnt.pool.ntp.org`
- **SNMP:** Not enabled.

## **5\. Login & Access**

- **Admin Username(s):** `ubnt`
- **Password Storage Method:** Stored in VaultWarden as "Ubiquiti Edge Router".

## **6\. Maintenance & Backup**

- **Configuration Backup:** No formal policy defined at this time; no ad-hoc backups performed.
- **Firmware Update Policy:** No formal policy defined at this time.

## **7\. Important Notes / Quirks**