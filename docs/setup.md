# Setup Guide

This guide covers everything needed to get the infrastructure running from scratch — from installing prerequisites to running your first `pulumi up`.

---

## Prerequisites

### Tools

- **Go 1.22+** — required to build the custom Pulumi providers
- **Pulumi CLI** — install from https://www.pulumi.com/docs/install/
- **Git**

### Accounts

- **Pulumi Cloud** — https://app.pulumi.com (free tier is sufficient)
- **Infisical** — https://app.infisical.com (for service-level secrets)
- **Cloudflare** — API token with `Zone:Read` and `DNS:Edit` permissions

### Host environment

- An **Unraid server** with Docker running
- Two bridged network interfaces: `br0` (primary LAN) and `br1` (container network)
- The following core services running before applying the apps stack:
  - **Nginx Proxy Manager** — `192.168.4.2`
  - **Technitium DNS** — `192.168.4.53`
  - **Authentik** — `192.168.4.55`

---

## 1. Clone the repository

```bash
git clone https://github.com/DCCoder90/home-net.git
cd home-net
```

---

## 2. Log in to Pulumi Cloud

```bash
pulumi login
```

This opens a browser to complete authentication. State will be stored in Pulumi Cloud going forward.

---

## 3. Create a stack

From the `pulumi/` directory:

```bash
cd pulumi
pulumi stack init <your-org>/home-net/dev
```

Replace `<your-org>` with your Pulumi Cloud organization name (shown in the top-left of the Pulumi Cloud UI).

---

## 4. Configure the stack

All non-secret values are already committed in `Pulumi.dev.yaml` and are picked up automatically. The only values you need to set in Pulumi config are the two **bootstrap credentials** required to authenticate with Infisical before anything else can run:

```bash
pulumi config set --secret infisicalClientId     <value>
pulumi config set --secret infisicalClientSecret <value>
```

> All other secrets (Cloudflare API token, provider credentials, admin email, etc.) are fetched automatically from Infisical at deploy time. Store them in your Infisical project under the `/secrets` path using the snake_case key names listed below.

### Infisical secret keys (`/secrets` folder)

These secrets are fetched from Infisical at runtime via the Universal Auth client ID/secret set in Pulumi config:

| Infisical key | Required | Description |
|---------------|----------|-------------|
| `cloudflare_api_token` | Yes | Cloudflare API token with `Zone:Read` and `DNS:Edit` permissions |
| `network_admin_email` | Yes | Email used for Let's Encrypt certificate requests |
| `technitium_admin_password` | Yes | Admin password for the Technitium DNS container |
| `ghcr_username` | Yes | GitHub Container Registry username |
| `ghcr_token` | Yes | GitHub Container Registry personal access token |
| `authentik_secret_key` | Yes | Authentik signing key (min 50 characters, generate with `openssl rand -hex 32`) |
| `authentik_postgresql_password` | Yes | PostgreSQL password for Authentik's database |
| `authentik_token` | **Phase 2 only** | Authentik API token — see two-phase workflow below |

Per-service runtime secrets (e.g. VPN credentials) and per-server SSH access keys (`/server_access` folder) are documented separately in the service YAML files.

### Additional Pulumi config values (beyond what's in `Pulumi.dev.yaml`)

In addition to `infisicalClientId` and `infisicalClientSecret`, set these in Pulumi config — they are not stored in Infisical:

```bash
# NPM Proxy Manager credentials
pulumi config set npmUsername <admin-email>
pulumi config set --secret npmPassword <admin-password>

# Technitium API token
pulumi config set --secret technitiumToken <api-token>

# Public-facing IP for external Cloudflare DNS records
pulumi config set publicFacingIp <your-public-ip>
```

### Two-phase Authentik bootstrap

Authentik cannot be configured via its API until it is running and initialized.
The deployment handles this with a two-phase workflow:

**Phase 1 — Deploy containers** (all secrets except `authentik_token` present):
```bash
pulumi up
```
This deploys all containers including the full Authentik stack (PostgreSQL, Redis, server, worker). Auth resources (groups, providers, policy bindings) are skipped. Services are deployed without auth protection.

**Initialize Authentik** (manual, one-time):
1. Navigate to `http://<authentik_ip>:9000` and complete initial setup
2. Create a service account and generate an API token
3. Add the token to Infisical as `authentik_token`

**Phase 2 — Configure auth** (`authentik_token` now present):
```bash
pulumi up
```
This run creates all Authentik auth resources (groups, providers, policy bindings) and enables auth protection on services that require it.

---

## 5. Build the custom providers

The two custom Pulumi providers (npmproxy and technitium) must be compiled before running any Pulumi commands. From the `pulumi/` directory:

```bash
bash scripts/build-providers.sh
```

This builds each provider binary and places it in `pulumi/bin/`, where `Pulumi.yaml` looks for local plugin overrides. The community Authentik provider (`OSMIT-GmbH/pulumi-authentik`) is fetched automatically by Pulumi from GitHub.

Re-run this script any time you modify a provider under `pulumi/providers/`.

---

## 6. First-run import (existing infrastructure only)

If this is a **fresh Pulumi state** importing pre-existing infrastructure, fill in the import IDs first:

```bash
# Run the helper script to auto-populate most IDs from running services:
export TOWER_HOST=<tower-ip> TOWER_SSH_USER=root \
       NPM_URL=http://192.168.4.2:81 NPM_USER=<user> NPM_PASS=<pass> \
       TECH_URL=http://192.168.4.53:5380 TECH_TOKEN=<token>
bash scripts/generate-imports.sh > ../config/tf-imports.yaml

# Then run with import mode enabled:
PULUMI_IMPORT_IDS_FILE=config/tf-imports.yaml pulumi up
```

After a successful import run, the `config/tf-imports.yaml` file can be left in place — empty-string entries are ignored and real IDs are only applied once (Pulumi skips import if the resource is already in state).

## 7. Preview and apply

```bash
# Dry-run — shows what will be created/modified/destroyed
pulumi preview

# Apply
pulumi up
```

---

## 8. Set up GitHub Actions (CI/CD)

The repository includes two workflows:

| Workflow | Trigger | Action |
|----------|---------|--------|
| `run-plan.yml` | PR against `main` | Runs `pulumi preview`, posts output as a PR comment |
| `run-apply.yml` | Push/merge to `main` | Runs `pulumi up --yes` |

Add the following secrets to your GitHub repository under **Settings → Secrets and variables → Actions**:

| Secret | Value |
|--------|-------|
| `PULUMI_ACCESS_TOKEN` | Create at https://app.pulumi.com/account/tokens |
| `PULUMI_STACK` | Fully-qualified stack name, e.g. `DCCoder90/home-net/dev` |

---

## 9. Adding a new service

1. Add a YAML file to `config/stacks/` (for a group of services) or `config/services/` (for a standalone service).
2. Follow the schema documented in:
   - [Stack Configuration Guide](config/stack-config.md)
   - [Service Configuration Guide](config/service-config.md)
3. Open a PR — the preview workflow will show exactly what Pulumi plans to create.
4. Merge to apply.
