# Wifi IP Pool

:::info
Current AP hardware does not support VLANs, as such there will only be a logical separation of networks rather than a VLAN separation.
:::

Current subnet `192.168.2.0/23` gives:

- Network Addresses: `192.168.2.0` to `192.168.3.255`
- Usable IP Range: `192.168.2.1` to `192.168.3.254`
- Total Usable IPs: 510

The best way to achieve this separation is by using **VLANs (Virtual LANs)**. Each category will be on its own VLAN, which means it will be its own broadcast domain and have its own subnet. You'll then typically assign different Wi-Fi SSIDs to each VLAN.

### **Proposed Subnet Split:**

1.  **Smart Home Devices:**
    - **Subnet:** `192.168.2.0/25`
    - **Network Address:** `192.168.2.0`
    - **Usable IP Range:** `192.168.2.1` - `192.168.2.126`
    - **Broadcast Address:** `192.168.2.127`
    - **Mask:** `255.255.255.128`
    - **Total Usable IPs:** 126
    - _Rationale:_ Smart home devices can be numerous. This provides a good number of addresses. These devices should often be isolated for security.
2.  **Utilities (phones, 3D printers, security system base station, etc.):**
    - **Subnet:** `192.168.2.128/26`
    - **Network Address:** `192.168.2.128`
    - **Usable IP Range:** `192.168.2.129` - `192.168.2.190`
    - **Broadcast Address:** `192.168.2.191`
    - **Mask:** `255.255.255.192`
    - **Total Usable IPs:** 62
    - _Rationale:_ For essential shared devices. This size should be sufficient for most households.
3.  **Free Range (Guests or uncategorized devices):**
    
    - **Subnet:** `192.168.2.192/26`
    - **Network Address:** `192.168.2.192`
    - **Usable IP Range:** `192.168.2.193` - `192.168.2.254`
    - **Broadcast Address:** `192.168.2.255`
    - **Mask:** `255.255.255.192`
    - **Total Usable IPs:** 62
    - _Rationale:_ For guest access or temporary devices. This keeps them isolated from main networks.
    
    _(The above three subnets fully utilize the_ `192.168.2.0/24` _portion of original /23 block)_
4.  **Kids Devices:**
    - **Subnet:** `192.168.3.0/25`
    - **Network Address:** `192.168.3.0`
    - **Usable IP Range:** `192.168.3.1` - `192.168.3.126`
    - **Broadcast Address:** `192.168.3.127`
    - **Mask:** `255.255.255.128`
    - **Total Usable IPs:** 126
    - _Rationale:_ Allows for multiple devices per child and future growth. Ideal for applying specific content filtering or time restrictions.
5.  **Adult Devices:**
    
    - **Subnet:** `192.168.3.128/25`
    - **Network Address:** `192.168.3.128`
    - **Usable IP Range:** `192.168.3.129` - `192.168.3.254`
    - **Broadcast Address:** `192.168.3.255`
    - **Mask:** `255.255.255.128`
    - **Total Usable IPs:** 126
    - _Rationale:_ Provides ample addresses for personal devices.
    
    _(The above two subnets fully utilize the_ `192.168.3.0/24` _portion of original /23 block)_

### **How to Implement:**

1.  **Router Configuration (e.g., EdgeRouter Pro):**
    - Define each of these subnets on your router.
    - Create a VLAN for each category (e.g., VLAN 10 for Kids, VLAN 20 for Smart Home, etc.).
    - Assign each new subnet to its corresponding VLAN interface on the router.
    - Configure a DHCP server on your router for each new subnet/VLAN to automatically assign IP addresses to devices connecting to that category.
2.  **Wi-Fi Access Points:**
    - Configure your Wi-Fi access points to broadcast multiple SSIDs (Wi-Fi names).
    - Map each SSID to the appropriate VLAN. For example:
        - `WiFi_Kids` -> Kids VLAN
        - `WiFi_SmartHome` -> Smart Home VLAN
        - `WiFi_Utilities` -> Utilities VLAN
        - `WiFi_Guest` -> Free Range VLAN
        - `WiFi_Adults` -> Adult Devices VLAN
    - Ensure your access points and any switches between the APs and the router are VLAN-aware and configured correctly (trunk ports for connections carrying multiple VLANs).
3.  **Firewall Rules:**
    - This is crucial. Once you have separate VLANs/subnets, you need to configure firewall rules on your EdgeRouter to control traffic _between_ these categories.
    - **Default:** Deny all inter-VLAN traffic unless explicitly allowed.
    - **Kids Devices:** Likely block access to Adult Devices, Utilities, and Smart Home admin interfaces. Allow internet access (possibly filtered).
    - **Smart Home Devices:** Isolate these as much as possible. They usually only need internet access (for cloud services) and perhaps limited access from specific devices on your trusted networks (e.g., your phone on the Adult Devices VLAN to control them). Block incoming connections from the internet.
    - **Utilities:** May need specific access rules. For example, phones might want to cast to a TV (which could be on another VLAN), or you might want to access a printer from multiple VLANs.
    - **Free Range (Guest):** Strictly isolate this VLAN. Allow internet access only. Block all access to any other internal VLAN/subnet.
    - **Adult Devices:** May have more permissive rules but should still be mindful of not unnecessarily exposing other more sensitive networks like Smart Home.

### **Benefits of this Approach:**

- **Security:** Isolate less secure devices (like some IoT/Smart Home devices or guest devices) from your more trusted devices.
- **Policy Enforcement:** Easily apply different rules (internet filtering, Quality of Service, access times) to different categories. For example, you can use the Technitium DNS groups we discussed earlier much more effectively by mapping these subnets to the appropriate DNS policy groups.
- **Network Management:** Better organization and easier troubleshooting.
- **Performance:** Reduces broadcast traffic within each segment.

Here's a table summarizing the IP allocation for each category:

|     |     |     |     |     |     |     |
| --- | --- | --- | --- | --- | --- | --- |
| **Category** | **Subnet (CIDR)** | **Network Address** | **Usable IP Range** | **Broadcast Address** | **Subnet Mask** | **Total Usable IPs** |
| Smart Home Devices | `192.168.2.0/25` | `192.168.2.0` | `192.168.2.1` - `192.168.2.126` | `192.168.2.127` | `255.255.255.128` | 126 |
| Utilities | `192.168.2.128/26` | `192.168.2.128` | `192.168.2.129` - `192.168.2.190` | `192.168.2.191` | `255.255.255.192` | 62  |
| Free Range (Guests) | `192.168.2.192/26` | `192.168.2.192` | `192.168.2.193` - `192.168.2.254` | `192.168.2.255` | `255.255.255.192` | 62  |
| Kids Devices | `192.168.3.0/25` | `192.168.3.0` | `192.168.3.1` - `192.168.3.126` | `192.168.3.127` | `255.255.255.128` | 126 |
| Adult Devices | `192.168.3.128/25` | `192.168.3.128` | `192.168.3.129` - `192.168.3.254` | `192.168.3.255` | `255.255.255.128` | 126 |