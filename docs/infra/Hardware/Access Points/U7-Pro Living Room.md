# U7-Pro (Living Room)

:::info
Last Updated 4/15/26
:::

## **1\. Device Overview**

- **Device Type:** Wireless Access Point
- **Manufacturer:** Ubiquiti
- **Model:** U7-Pro
- **Primary Role:** Provides WiFi network access for wireless devices in the living room area.
- **Operation Mode:** Access Point
- **Management:** Managed via UniFi Network Controller on [Cloud Gateway Max](../Routers/UCG-Max.md)

## **2\. Physical Information**

- **Physical Location:** Living Room
- **Power Method:** PoE — powered directly from [USW-Pro-Max-16-PoE](../Switches/USW-Pro-Max-16-PoE.md)
- **Uplink Connection:** USW-Pro-Max-16-PoE _(port to be added)_

## **3\. Network Configuration**

- **IP Address:** `192.168.1.130`
- **IP Assignment:** DHCP Reservation on [Cloud Gateway Max](../Routers/UCG-Max.md)
- **MAC Address:** `8c:ed:e1:84:f9:30`
- **Subnet:** `192.168.1.0/24`
- **Default Gateway:** `192.168.1.1`
- **DNS Servers:** `192.168.4.53` (Technitium DNS)

## **4\. Wireless Configuration**

- **Main WiFi Network:**
    - **SSID (Network Name):** Stored in VaultWarden under "ui.com"
    - **Password:** Stored in VaultWarden under "ui.com"
    - **Security Type:** WPA3/WPA2 (managed by UniFi)
    - **SSID Broadcast:** Enabled
- **Guest Network:** Not enabled / Not in use.

## **5\. Login & Access**

- **Admin Password Storage:** Stored in VaultWarden under "ui.com"

## **6\. Maintenance & Backup**

- **Configuration Backup:** Managed by UniFi Network Controller on [Cloud Gateway Max](../Routers/UCG-Max.md).

## **7\. Important Notes / Quirks**
