# Home-Net

This repository contains configurations for managing various services and infrastructure components within my home network.

## Purpose

I want an easy way to automate and keep track of my local home lab.  It would be nice to also have a way to easily and quickly recreate it in case of a catastrophic failure (flood, house fire, etc etc)

## Modules

Reusable Terraform modules are located in the `modules/` directory. These modules are designed to be configurable and handle specific types of deployments.

## Services

The `services/` directory contains specific Terraform configurations that deploy various applications and services within the home network. These configurations make use of the reusable modules defined in the `modules/` directory. Each subdirectory within `services/` typically represents a distinct service or a group of related services.

## Prerequisites

*   **Target Environment:** This setup is primarily designed to run on an **Unraid server**.
*   **Terraform Agent:** A Terraform agent [must be running on the Unraid server](https://developer.hashicorp.com/terraform/cloud-docs/agents/agents#run-an-agent-with-docker).
    *   This agent requires **direct access to the Docker socket** (typically `/var/run/docker.sock`) to manage Docker resources. Ensure the user running the agent has the necessary permissions.
*   **Terraform CLI:** Terraform (version compatible with the configurations, e.g., v1.12.2 or later) installed.
*   **Docker:** Docker must be configured and running on the Unraid server.

## Usage

1.  **Clone the repository:**
    ```bash
    git clone <repository-url>
    cd home-net
    ```
2.  **Configure Backend:**
    Ensure your Terraform backend (e.g., HCP Terraform, S3) is configured in your root `main.tf` file.
3.  **Initialize Terraform:**
    ```bash
    terraform init
    ```
4.  **Review:**
    ```bash
    terraform plan
    ```

Any applies should be done using the [run-apply](./.github/workflows/run-apply.yml) action.

## How is this run?

This is run on a self-hosted agent on the Unraid Server.  This is invoked using HCP Terraform, the state is managed by HCP Terraform and I have the following variables set:

|Variable Name|Type|Sensitive?|
|---|---|---|
|cloudflare_api_token|terraform|Y|
|network_admin_email|terraform|Y|
|nginx_proxy_address|terraform|N|
|nginx_proxy_pass|terraform|Y|
|nginx_proxy_user|terraform|Y|
|public_facing_ip|terraform|Y|
|technitium_api_token|terraform|Y|
|technitium_host|terraform|N|
|AUTHENTIK_INSECURE|env|N|
|AUTHENTIK_TOKEN|env|Y|
|AUTHENTIK_URL|env|N|
|vpn_pass|terraform|Y|
|vpn_user|terraform|Y|


You may be wondering why some things such as "public facing API" and "network admin email" are set to sensitive.  Well, I honestly I don't want the world knowing those details.  Aside from that, there's no reason for them to be marked sensitive.


### Cloudflare

The Cloudflare API token must have these permissions:
- Zone:Read
- DNS:Edit

In addition it must also be able to access any zones that are wished to be used.