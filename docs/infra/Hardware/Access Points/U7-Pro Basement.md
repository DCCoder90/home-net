# U7-Pro (Basement)

:::info
Last Updated 4/15/26
:::

## **1\. Device Overview**

- **Device Type:** Wireless Access Point
- **Manufacturer:** Ubiquiti
- **Model:** U7-Pro
- **Primary Role:** Provides WiFi network access for wireless devices in the basement area.
- **Operation Mode:** Access Point
- **Management:** Managed via UniFi Network Controller on [Cloud Gateway Max](../Routers/UCG-Max.md)

## **2\. Physical Information**

- **Physical Location:** Basement Rack
- **Power Method:** PoE — powered directly from [USW-Pro-Max-16-PoE](../Switches/USW-Pro-Max-16-PoE.md) Port 9
- **Uplink Connection:** USW-Pro-Max-16-PoE Port 9

## **3\. Network Configuration**

- **IP Address:** `192.168.1.214`
- **IP Assignment:** DHCP Reservation on [Cloud Gateway Max](../Routers/UCG-Max.md)
- **MAC Address:** `8c:ed:e1:84:fa:14`
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
- **Firmware Update Policy:** No formal policy defined at this time.

## **7\. Important Notes / Quirks**
