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

The non-secret values are already committed in `Pulumi.dev.yaml` and are picked up automatically. You only need to set the secrets:

```bash
pulumi config set --secret infisicalClientId       <value>
pulumi config set --secret infisicalClientSecret   <value>
pulumi config set --secret nginxProxyUser          <value>
pulumi config set --secret nginxProxyPass          <value>
pulumi config set --secret cloudflareApiToken      <value>
pulumi config set --secret technitiumApiToken      <value>
pulumi config set --secret technitiumAdminPassword <value>
pulumi config set --secret vpnUser                 <value>
pulumi config set --secret vpnPass                 <value>
pulumi config set --secret ghcrToken               <value>
pulumi config set --secret authentikToken          <value>
pulumi config set --secret networkAdminEmail       <value>
pulumi config set --secret publicFacingIp          <value>
```

> Secrets are encrypted by Pulumi Cloud and never stored in plain text. Running `pulumi config` will show them as `[secret]`.

---

## 5. Build the custom providers

The three custom Pulumi providers (npmproxy, technitium, authentik) must be compiled and installed before running any Pulumi commands:

```bash
bash scripts/build-providers.sh
```

This builds each provider from source and installs the binaries to `~/.pulumi/plugins/`. Re-run this script any time you modify a provider under `pulumi/providers/`.

---

## 6. Preview and apply

```bash
# Dry-run — shows what will be created/modified/destroyed
pulumi preview

# Apply
pulumi up
```

---

## 7. Set up GitHub Actions (CI/CD)

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

## Adding a new service

1. Add a YAML file to `config/stacks/` (for a group of services) or `config/services/` (for a standalone service).
2. Follow the schema documented in:
   - [Stack Configuration Guide](config/stack-config.md)
   - [Service Configuration Guide](config/service-config.md)
3. Open a PR — the preview workflow will show exactly what Pulumi plans to create.
4. Merge to apply.
