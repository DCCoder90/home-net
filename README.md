# Home-Net

This repository contains Terraform configurations for managing various services and infrastructure components within my home network.

## Modules

Reusable Terraform modules are located in the `modules/` directory. These modules are designed to be configurable and handle specific types of deployments.

## Services

The `services/` directory contains specific Terraform configurations that deploy various applications and services within the home network. These configurations make use of the reusable modules defined in the `modules/` directory. Each subdirectory within `services/` typically represents a distinct service or a group of related services.

## Prerequisites

*   Terraform (version compatible with the configurations, e.g., v1.12.2 or later)
*   Docker installed and running on the target host(s) if using the Docker module.

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