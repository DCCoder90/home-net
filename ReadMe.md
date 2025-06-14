# Home-Net

This repository contains Terraform configurations for managing various services and infrastructure components within my home network.

Quick note:

Cloudflare API token must have these permissions:
- Zone:Read
- DNS:Edit


## Modules

Reusable Terraform modules are located in the `modules/` directory. These modules are designed to be configurable and handle specific types of deployments.

## Services

The `services/` directory contains specific Terraform configurations that deploy various applications and services within the home network. These configurations make use of the reusable modules defined in the `modules/` directory. Each subdirectory within `services/` typically represents a distinct service or a group of related services.

## Prerequisites

*   **Target Environment:** This setup is primarily designed to run on an **Unraid server**.
*   **Terraform Agent:** A Terraform agent (e.g., HCP Terraform self-hosted agent, GitLab Runner, GitHub Actions self-hosted runner, etc.) must be running on the Unraid server.
    *   This agent requires **direct access to the Docker socket** (typically `/var/run/docker.sock`) to manage Docker resources. Ensure the user running the agent has the necessary permissions.
*   **Terraform CLI:** Terraform (version compatible with the configurations, e.g., v1.12.2 or later) installed on the machine where `terraform plan/apply` commands are initiated or on the agent itself.
*   **Docker:** Docker must be installed and running on the Unraid server.

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
4.  **Review and Apply:**
    ```bash
    terraform plan
    terraform apply
    ```