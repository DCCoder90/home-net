# USW-Pro-Max-16-PoE

:::info
Last Updated 4/15/26
:::

## **1\. Device Overview**

- **Hostname/Device Name:** _(To be added)_
- **Model:** Ubiquiti USW-Pro-Max-16-PoE (180W)
- **Role:** Primary Access Layer Switch
- **Management Access:**
    - **IP Address:** `192.168.1.253`
    - **Managed via:** UniFi Network Controller on [Cloud Gateway Max](../Routers/UCG-Max.md)

## **2\. Physical Information**

- **Physical Location:** Basement Rack
- **Serial Number:** [To be added later]
- **Power Backup (UPS):**
    - UPS Model: [To be added later]
- **Replaces:** Linksys LGS124 24-Port Gigabit Unmanaged & Ubiquiti EdgeSwitch 10X

## **3\. Port Layout**

| Port | Type | Connection |
|------|------|------------|
| 1 | RJ45 (PoE) | [Cloud Gateway Max](../Routers/UCG-Max.md) |
| 2 | RJ45 (PoE) | Not Connected |
| 3 | RJ45 (PoE) | Tower (MAC: `...72:bc`) |
| 4 | RJ45 (PoE) | Not Connected |
| 5 | RJ45 (PoE) | PlayStation 5 |
| 6 | RJ45 (PoE) | Not Connected |
| 7 | RJ45 (PoE) | Not Connected |
| 8 | RJ45 (PoE) | Eufy Base Station |
| 9 | RJ45 (PoE) | [U7-Pro (Basement)](../Access%20Points/U7-Pro%20Basement.md) |
| 10 | RJ45 (PoE) | Not Connected |
| 11 | RJ45 (PoE) | Not Connected |
| 12 | RJ45 (PoE) | Not Connected |
| 13 | RJ45 (PoE) | Not Connected |
| 14 | RJ45 (PoE) | Not Connected |
| 15 | RJ45 (PoE) | Not Connected |
| 16 | RJ45 (PoE) | Not Connected |
| 17 | SFP+ | Not Connected |
| 18 | SFP+ | Not Connected |

## **4\. VLAN Configuration**

- **VLANs:** Not currently configured.

## **5\. Login & Access**

- **Admin Password Storage:** Stored in VaultWarden under "ui.com"

## **6\. Maintenance & Backup**

- **Configuration Backup:** No formal policy defined at this time.

## **7\. Important Notes / Quirks**
