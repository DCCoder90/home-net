# Tower

:::info
Last Updated 7/8/2025
:::

## **1\. Overview**

- **Hostname:** tower
- **Primary Role:** Main Unraid NAS and Application Server for Homelab.
- **Operating System:** Unraid OS
- **OS Version:** `7.1.2` (Build Date: `2025-05-11`)
- **Default Management Access:** WebUI via
- `http://tower/`
- `http://192.168.1.41/` (eth0)
- `http://192.168.4.41/` (eth1)

## **2\. Hardware Specifications**

- **Motherboard:**
    - **Manufacturer:** Gigabyte Technology Co., Ltd.
    - **Model:** A620M GAMING X AX
    - **Version:** x.x
    - **BIOS Manufacturer:** American Megatrends International, LLC.
    - **BIOS Version:** FA
    - **BIOS Date:** Thu 22 Feb 2024
- **CPU:**
    - **Model:** AMD Ryzen 5 8500G w/ Radeon 740M Graphics
    - **Base Clock Speed:** 3550 MHz
- **Integrated GPU (iGPU):** Radeon 740M Graphics (Part of Ryzen 5 8500G)
    - **Usage:** Not actively used.
- **Discrete GPU (dGPU):**
    - **Model:** NVIDIA GeForce GTX 1050 Ti
    - **Usage:** AI tasks and video/audio transcoding.
- **Memory (RAM):**
    - **Capacity:** 32 GiB
    - **Type:** DDR5
- **Storage:**
    - **Boot Device:**
        - **Model:** FlashUSB_3.2.1_FD
        - **Capacity:** 31 GB
        - **Device ID:** `sda`
    - **Parity Drive(s):**
        - `sdc`: ST4000VN006-3CW104_ZW63EE0A - 4 TB
    - **Array Data Drives:**
        - `sde`: ST4000VN008-2DR166_ZDHBCGYM - 4 TB
        - `sdb`: ST4000VN006-3CW104_ZW63EFFH - 4 TB
        - `sdd`: ST4000VN006-3CW104_ZW63AVXW - 4 TB
    - **Cache Drives / Pools:** None configured.
    - **Unassigned Devices:** None currently configured.

## **3\. Network Configuration**

- **Interface** `eth0` **(Wired):**
    - **MAC Address:** `10:FF:E0:24:7F:EE`
    - **IP Address:** `192.168.1.41`
    - **Purpose:** Access to the rest of the home network.
- **Interface** `eth1` **(Wired):**
    - **MAC Address:** `68:05:CA:BE:72:BC`
    - **IP Address:** `192.168.4.41`.
    - **Purpose:** Private network, physical connection to eth4 on [EdgeRouter Pro](../Routers/EdgeRouter%20Pro.md)
- **Interface** `wlan0` **(Wireless):**
    - **MAC Address:** _(Not currently on hand)_
    - **IP Address:** `192.168.2.208`
    - **Purpose:** Access to the rest of the home network, intended as redundant for `eth0`.
- **VPN / Overlay Network:**
    - **Tailscale:** [Tailscale](../../Services/External%20Services/Tailscale.md) Enabled on this device.
        - **Tailscale IP Address:** _(You can find this in the Tailscale admin console or on the device:_ `tailscale ip -4`_)_
        - **Purpose:** Provides remote access to the server and its services over the [Tailscale](../../Services/External%20Services/Tailscale.md) network.
- **Bonding/Link Aggregation:** Not configured.
- **DNS Servers Used by Tower:** `1.1.1.1`, `8.8.8.8`

## **4\. Key Roles & Responsibilities**

- Primary Network Attached Storage (NAS)
- Docker Application Host (via Unraid Apps)
- AI Processing tasks
- Video/Audio Transcoding services

## **5\. Physical Information**

- **Physical Location:** Beside the rack located under the basement stairs.
- **Power Backup (UPS):**
    - **Type:** Connected to a dedicated UPS.
    - **Rating:** 1500VA
    - **Management:** Configured via "Network UPS Tools" (NUT) from Unraid Apps, connected via USB cable for graceful shutdown.

## **6\. Login & Access**

- **WebUI Username:** `root`
- **WebUI Password:** Stored in Vaultwarden under "Unraid Server Root".
- **SSH Access:** Not enabled.
- **Unraid Connect:** Not enabled.

## **7\. Maintenance & Notes**

- **BIOS Settings:** Default BIOS settings are in use.
- **Regular Checks:**
    - **Parity Check Schedule:** Weekly, every Sunday at 12:00 AM (Midnight).
- **Important Notes:**
    - The address range of `192.168.4.0/23` is reserved by the [EdgeRouter Pro](../Routers/EdgeRouter%20Pro.md) for allocation by this server for running services.