# Home

:::info
Last Updated 6/02/25
:::

## **Purpose of This Space**

This serves as the central repository for all documentation related to my personal homelab environment. The goal is to maintain a clear, up-to-date record of its architecture, components, configurations, and the services it provides.

This documentation is intended to:

- Aid in troubleshooting and problem-solving.
- Provide a reference for how systems are configured and interconnected.
- Track changes and new additions to the lab.
- Serve as a personal knowledge base for the various technologies implemented.

## **Scope of Documentation**

Within this space, you will find detailed information covering various aspects of the homelab, including:

- **Networking Equipment:** Configuration and details for routers, switches, wireless access points, and other network infrastructure.
- **Server Hardware:** Specifications, configurations, and roles of physical servers, including the primary Unraid server ("tower").
- **Virtualization & Containerization:** Details on virtual machines and Dockerized services.
- **Storage Systems:** Information on NAS devices, storage pools, and data management.
- **Services:** Dedicated pages for each service running within the homelab, outlining their purpose, setup, access details, and dependencies. This includes media servers, automation tools, security applications, and more.
- **Configuration Details:** Specific settings, important parameters, and setup notes for various hardware and software components.

**Note on Scope:** This documentation focuses specifically on the homelab infrastructure and its dedicated services. Standard consumer-grade end-user devices such as personal computers, tablets, smartphones, smart TVs, and similar personal electronics are considered outside the scope of this documentation, unless they directly host or critically interact with a documented homelab service.

## **General Homelab Architecture Overview**

My homelab primarily utilizes an Unraid server, named **"tower"**, for hosting the majority of services. Most applications are installed and managed as Docker containers, often through Unraid's "Apps" interface (Community Applications).

- **Default Deployment:** Unless specified otherwise on an individual service or hardware page, assume components are part of or managed by the "tower" Unraid server, with services running as Docker containers.
- **Exceptions:** Any services or hardware running on dedicated Virtual Machines, different physical hardware (like the EdgeRouter or Synology NAS), or external platforms will have their specific hosting environment and setup detailed on their respective documentation pages.

## Navigating This Space

This documentation is broadly divided into two main sections: **Home**, **Hardware** and **Services**.

- [Home](Home.md) - The page you are currently viewing.
    - [Future Enhancements](Home/Future%20Enhancements.md) - This page contains a wish list of future enhancements to the homelab or it’s network that have yet to be completed.
    - [Documentation](Home/Documentation.md) - This page contains important notes about the formatting of this documentation.
    - [IP Pool (tower reserved)](Home/IP%20Pool%20%28tower%20reserved%29.md) - This page contains important notes about IP allocation for services.

- [Hardware](Hardware.md) **:** This section contains documentation for the physical components of the homelab. It is further organized into subcategories such as:
    - [Servers](Hardware/Servers.md)
    - [Access Points](Hardware/Access%20Points.md)
    - [Switches](Hardware/Switches.md)
    - [Routers](Hardware/Routers.md)
    - [Other Hardware](Hardware/Other%20Hardware.md)
- [Services](Services.md) **:** This section details the software applications and services running within the homelab. It is further organized into logical subcategories like:
    - [Networking](Services/Networking.md)
    - [Media](Services/Media.md)
    - [Database](Services/Database.md)
    - [External Services](Services/External%20Services.md)
    - [Other Services](Services/Other%20Services.md)

Individual pages are typically named after the specific device or service.

This is a living document and will be updated as the homelab evolves.