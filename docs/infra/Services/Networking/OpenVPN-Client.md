# OpenVPN-Client

:::info
Last Updated 5/27/25
:::

## **1\. Service Overview**

- **Service Name:** OpenVPN-Client (User-defined name for container)
- **Container Maintainer/Image Source:** `ich777`
- **Primary Purpose:** To establish an outbound VPN connection using [ProtonVPN](../External%20Services/ProtonVPN.md) . The primary use is to serve as a [Tailscale](../External%20Services/Tailscale.md) exit node, allowing [Tailscale](../External%20Services/Tailscale.md) clients to route their internet traffic through [ProtonVPN](../External%20Services/ProtonVPN.md) for privacy and to bypass region restrictions.
- **Status:** Production

## **2\. Host & Docker Details**

- **Running On:** [Tower](../../Hardware/Servers/Tower.md)
- **Unraid App Name (if applicable):** OpenVPN-Client _(This is the container name in Unraid's Docker management. If from Community Applications, the specific CA App name could be added.)_
- **Docker Image:** `ghcr.io/ich777/openvpn-client`
- **Docker Hub Registry:** https://hub.docker.com/r/ich777/openvpn-client/
- **Network Type:** `Bridge`
- **TUN Device Passthrough:** Host path `/dev/net/tun` is mapped to the container.
- **Key Docker Run Parameters / Unraid Extra Parameters:**
    - `--cap-add=NET_ADMIN` (Required for network manipulation)
    - `--dns=8.8.8.8` (Sets container's internal DNS resolver)
    - `--sysctl net.ipv6.conf.all.disable_ipv6=1` (Disables IPv6 within the container)
    - `--restart=unless-stopped` (Container restart policy)

## **3\. Configuration & Data**

- **Volume Mounts (Data Persistence):**
    - **Container Path:** `/vpn`
    - **Host Path (Unraid):** `/mnt/user/appdata/openvpn-client`
- **VPN Provider:** NordVPN
    - **NordVPN Configuration Method:** Configuration is managed through two files located in the host path mapped to `/vpn` (i.e., within `/mnt/user/appdata/openvpn-client/`, typically in an `openvpn` subdirectory or directly in the root of the mapped volume, depending on image specifics):
        - `vpn.ovpn`: The OpenVPN configuration file provided by NordVPN for the desired server/settings. See [vpn.ovpn](OpenVPN-Client/vpn.ovpn.md)
        - `vpn.auth`: A file containing NordVPN username and password, typically two lines (username on the first, service password on the second).
- **Key Environment Variables (Unraid Template or Docker Run):**
    - `PING_INTERVAL=30`
- **Backup Strategy:** The configuration files (`vpn.ovpn`, `vpn.auth`) located in `/mnt/user/appdata/openvpn-client/` are **not currently backed up**.

## **4\. Networking**

- **Container IP Address:** As the network type is "Bridge," the container receives its own IP address on a private Docker bridge network (e.g., `172.17.0.x`). Outbound traffic from the container (including the VPN tunnelled traffic) is Network Address Translated (NATted) by Docker to appear as if originating from the Unraid host's ("tower") IP address. The specific internal Docker IP of the container is generally not needed for its function as an exit node.
- **DNS Resolver (Container Internal):** `8.8.8.8` (as per extra parameters)
- **IPv6:** Disabled within the container.
- **Purpose/Users of the VPN Connection:** Users on the Tailscale network can select this "vpn" device as an exit node to route their internet traffic through the NordVPN connection established by this container. This is used for enhanced privacy and bypassing region-based content restrictions.

## **5\. Tailscale Integration**

- **Tailscale Status:** Enabled.
- **Tailscale Hostname:** `vpn` (This is the name of this node on your Tailscale network)
- **Tailscale Role:** Configured as an **Exit Node**.

## **6\. Management & Troubleshooting**

- **Update Process:** Via Unraid Apps tab -> Check for Updates if installed as a CA app, or by manually pulling the new Docker image (`ghcr.io/ich777/openvpn-client:latest`) and recreating the container if managed manually.
- **Restart Procedure:** Via Unraid Docker UI (Restart container).
- **Log Access:**
    - Via Unraid Docker UI (Logs button for the container).
- **Checking Connection Status:** For clients using this as a Tailscale exit node, their public IP address can be checked using an external website like `whatismyip.com` to verify it matches the NordVPN server's IP.

## **7\. Dependencies**

- **Unraid Server "tower":** For Docker hosting and device passthrough.
- **Internet Connectivity:** Required for [ProtonVPN](../External%20Services/ProtonVPN.md) and [Tailscale](../External%20Services/Tailscale.md) .
- **NordVPN Service:** Relies on an active [ProtonVPN](../External%20Services/ProtonVPN.md) subscription and their server infrastructure. Requires correctly configured `vpn.ovpn` and `vpn.auth` files.
- **Tailscale Service:** Relies on Tailscale's coordination servers (unless using Headscale).

## **8\. Useful Links & Notes**

- **Docker Hub:** https://hub.docker.com/r/ich777/openvpn-client/
- **GitHub (ich777):** https://github.com/ich777/docker-openvpn-client (Primary source for `ghcr.io` image)