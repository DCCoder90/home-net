# Alerts

:::info
Last Updated 06/02/25
:::

### **1\. Introduction**

This document outlines the alerts configured for the Unraid server. The purpose of these alerts is to proactively notify administrators of potential issues or important events related to server health and resource utilization. Early detection allows for timely intervention, ensuring server stability and performance.

### **2\. Alerting System Overview**

- **Alert Evaluation & Setup Location:** The alerts detailed below are configured and managed within **Grafana's alerting system**.
- **Notification Channel:** All triggered alerts send notifications via a **Discord webhook**.

### **3\. Unraid Alert Definitions**

- **Metric Source:** All metrics for these alerts are primarily sourced from the **Unraid node exporter**.

#### **3.1.** [UnraidHighCpuUsage](Alerts/UnraidHighCpuUsage.md)

- **Purpose:** To notify when the overall CPU utilization on the Unraid server is excessively high for a sustained period. Prolonged high CPU usage can indicate a runaway process, insufficient processing power for the current workload, or other performance bottlenecks that may impact server responsiveness and the performance of hosted applications.
- **Trigger Condition:** CPU utilization is greater than **90%**.
- **Evaluation Details:**
    - **Evaluation Interval:** The alert condition is checked every **1 minute**.
    - **Pending Period (Time to Fire):** The condition must persist for **5 minutes** before an alert is triggered and a notification is sent.
    - **Keep Firing:** 0s (Notification sent once when alert becomes active).
- **Severity:** Warning
- **Notification:** Sends a message to Discord.
- **Metric Source:** Unraid node exporter.
- **Setup Location:** Grafana Alerting.

#### **3.2.** [Disk1SpaceWarning](Alerts/Disk1SpaceWarning.md)

- **Purpose:** To warn when the filesystem associated with Unraid's "Disk 1" is nearing its capacity. Running out of disk space can lead to data loss, application crashes, and system instability. This alert provides an early warning to allow for cleanup or capacity expansion.
- **Trigger Condition:** Used disk space on Disk 1 is greater than **95%**.
- **Evaluation Details:**
    - **Evaluation Interval:** The alert condition is checked every **1 minute**.
    - **Pending Period (Time to Fire):** The condition must persist for **1 minute** before an alert is triggered and a notification is sent.
    - **Keep Firing:** 0s (Notification sent once when alert becomes active).
- **Severity:** Warning
- **Notification:** Sends a message to Discord.
- **Metric Source:** Unraid node exporter.
- **Setup Location:** Grafana Alerting.

#### **3.3.** [Disk2SpaceWarning](Alerts/Disk2SpaceWarning.md)

- **Purpose:** To warn when the filesystem associated with Unraid's "Disk 2" is nearing its capacity. Similar to Disk 1, this aims to prevent issues arising from a full disk.
- **Trigger Condition:** Used disk space on Disk 2 is greater than **95%**.
- **Evaluation Details:**
    - **Evaluation Interval:** The alert condition is checked every **1 minute**.
    - **Pending Period (Time to Fire):** The condition must persist for **1 minute** before an alert is triggered and a notification is sent.
    - **Keep Firing:** 0s (Notification sent once when alert becomes active).
- **Severity:** Warning
- **Notification:** Sends a message to Discord.
- **Metric Source:** Unraid node exporter.
- **Setup Location:** Grafana Alerting.

#### **3.4.** [Disk3SpaceWarning](Alerts/Disk3SpaceWarning.md)

- **Purpose:** To warn when the filesystem associated with Unraid's "Disk 3" is nearing its capacity. Similar to Disk 1 and Disk 2, this aims to prevent issues arising from a full disk.
- **Trigger Condition:** Used disk space on Disk 3 is greater than **95%**.
- **Evaluation Details:**
    - **Evaluation Interval:** The alert condition is checked every **1 minute**.
    - **Pending Period (Time to Fire):** The condition must persist for **1 minute** before an alert is triggered and a notification is sent.
    - **Keep Firing:** 0s (Notification sent once when alert becomes active).
- **Severity:** Warning
- **Notification:** Sends a message to Discord.
- **Metric Source:** Unraid node exporter.
- **Setup Location:** Grafana Alerting.

### **4\. General Alert Management**

- **Viewing Alerts:** Active alerts can typically be viewed within the Grafana UI, under the "Alerting" section.
- **Responding to Alerts:** When an alert notification is received, investigate the condition promptly.