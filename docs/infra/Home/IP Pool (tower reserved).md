# IP Pool (tower reserved)

## **1\. Overview & Purpose**

This page documents the dedicated IP address pool `192.168.4.0/23`. This network segment is primarily reserved for use by the Unraid server [Tower](../Hardware/Servers/Tower.md) to assign IP addresses to Docker containers, Virtual Machines (VMs), or other services running directly on it and connected to its `br1` bridge network interface.

The primary goal of this dedicated network is to provide isolated and organized IP addressing for services hosted on [Tower](../Hardware/Servers/Tower.md)

## **2\. Network Details**

- **Network Address:** `192.168.4.0/23`
- **Subnet Mask:** `255.255.254.0`
- **Full IP Address Range:** `192.168.4.0` - `192.168.5.255`
- **Usable IP Address Range:** `192.168.4.1` - `192.168.5.254`
- **Broadcast Address:** `192.168.5.255`
- **Total Usable IPs:** 510

## **3\. Allocation & Management**

- **Reserved By:** The `192.168.4.0/23` IP range is reserved on the [EdgeRouter Pro](../Hardware/Routers/EdgeRouter%20Pro.md) . This means the EdgeRouter's DHCP servers for other LAN segments will not assign addresses from this pool.
- **Primary User:** Unraid Server **"tower"**.
- **Tower's Network Interface for this Segment:** `br1` (a bridge interface on "tower").
    - The physical interface `eth1` on "tower" is part of this `br1` bridge.
- **Tower's IP Address on** `br1`**:** `192.168.4.41` (Falls within the "Networking" sub-range below).

## **4\. IP Allocation Plan within** `192.168.4.0/23`

The following plan outlines the general allocation of IP address sub-ranges within this network for different categories of services. This is intended to help organize IP assignments.

|     |     |     |     |
| --- | --- | --- | --- |
| **Category** | **Start IP Address** | **Stop IP Address** | **Notes** |
| **Networking** | `192.168.4.1` | `192.168.4.50` | e.g., Proxies, DNS, network utilities. Tower's `br1` IP (`192.168.4.41`) is in this range. |
| **Infrastructure** | `192.168.4.51` | `192.168.4.100` | e.g., Authentication, management tools. |
| **Databases** | `192.168.4.101` | `192.168.4.150` |     |
| **Development** | `192.168.4.151` | `192.168.4.255` | e.g., Test instances, development tools. |
| **Media** | `192.168.5.1` | `192.168.5.50` | e.g., Plex, \*arr stack, Jellyfin. |
| **Reserved** | `192.168.5.51` | `192.168.5.254` | For future use or dynamic assignments if any. |

## **5\. Usage Notes**

- **Intended Use:** This IP range is designated for services running on "tower" that require their own IP addresses separate from "tower's" primary management IP (`192.168.1.41`). This is typically used for Docker containers (often with custom Docker networks attached to `br1`) or Virtual Machines.
- **IP Assignment to Services:**
    - IPs within the allocated sub-ranges are generally expected to be **statically assigned** to Docker containers or VMs within their respective configurations on "tower".
- **Gateway for this Network:**
    - _(Is_ `192.168.4.41` _(tower's br1 IP) the gateway for devices/containers on this_ `192.168.4.0/23` _network? Or does the EdgeRouter-Pro-8-Port have an interface/IP on this segment that acts as the gateway? This needs to be defined for outbound connectivity from these services.)_
- **Routing:** The EdgeRouter-Pro-8-Port is aware of how to route traffic to the `192.168.4.0/23` network, via "tower's" `eth1` MAC address.

## **6\. Services Utilizing This Network**

See [Services](../Services.md) for current IP allocations