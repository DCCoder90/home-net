# ARRR Stack Setup
This folder contains the Terraform configuration files responsible for deploying and managing the "Arrr" stack of services. This setup leverages Docker containers, secured by an Nginx reverse proxy with SSL, and integrated with Authentik for authentication. Local DNS entries are also configured to simplify access from a local network.

## Services Deployed
The following services are part of this deployment:

Prowlarr
Sonarr
Radarr
Lidarr
Readarr

Each service runs as a Docker container with PUID=1000, PGID=1000, and TZ=America/Chicago environment variables set.

## Terraform Modules Used
This setup is orchestrated using Terraform and relies on the following local modules:

**Authentik Proxy Authentication:** Configured via ../../modules/proxy_auth
**Nginx Proxy, SSL, and Local DNS:** Managed by ../../modules/dns
**Docker Containers:** Defined in ../../modules/docker